package service

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/app/repo"
	"github.com/1Panel-dev/1Panel/agent/buserr"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
	"github.com/1Panel-dev/1Panel/agent/utils/copier"
	"github.com/1Panel-dev/1Panel/agent/utils/ssh"
)

const defaultSSHKeyName = "id_ed25519"

var (
	sshHostAliasPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)
	sshHostUserPattern  = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
)

type ISSHHostService interface {
	Search(req dto.SearchWithPage) (int64, interface{}, error)
	Create(req dto.SSHHostOperate) (*dto.SSHHostInfo, error)
	Delete(req dto.ForceDelete) error
	Test(req dto.SSHHostTest) bool
}

type SSHHostService struct{}

func NewISSHHostService() ISSHHostService {
	return &SSHHostService{}
}

func (s *SSHHostService) Search(req dto.SearchWithPage) (int64, interface{}, error) {
	total, hosts, err := sshHostRepo.Page(req.Page, req.PageSize, sshHostRepo.WithByInfo(req.Info))
	if err != nil {
		return 0, nil, err
	}
	var data []dto.SSHHostInfo
	for _, item := range hosts {
		var info dto.SSHHostInfo
		_ = copier.Copy(&info, &item)
		data = append(data, info)
	}
	return total, data, nil
}

func (s *SSHHostService) Create(req dto.SSHHostOperate) (*dto.SSHHostInfo, error) {
	if err := validateSSHHostInput(req.Alias, req.HostName, req.User); err != nil {
		return nil, err
	}
	password, err := decodeSSHHostPassword(req.Password)
	if err != nil {
		return nil, err
	}
	if len(password) == 0 {
		return nil, buserr.New("ErrSSHHostPasswordRequired")
	}

	if existing, _ := sshHostRepo.Get(sshHostRepo.WithByAlias(req.Alias)); existing.ID != 0 {
		return nil, buserr.New("ErrRecordExist")
	}

	configPath, err := sshConfigPath()
	if err != nil {
		return nil, err
	}
	configContent, err := readSSHConfigFile(configPath)
	if err != nil {
		return nil, err
	}
	if ssh.HasSSHConfigHost(configContent, req.Alias) {
		return nil, buserr.New("ErrRecordExist")
	}

	if err := ensureDefaultSSHKey(); err != nil {
		return nil, err
	}
	publicKey, privateKey, err := loadDefaultSSHKeyPair()
	if err != nil {
		return nil, err
	}

	port := req.Port
	if port == 0 {
		port = 22
	}
	if err := ssh.CopyPublicKeyWithPassword(req.HostName, req.User, password, port, publicKey); err != nil {
		return nil, fmt.Errorf("ssh-copy-id failed: %w", err)
	}

	updatedConfig, err := ssh.AddSSHConfigHost(configContent, req.Alias, req.HostName, req.User)
	if err != nil {
		return nil, err
	}
	if err := writeSSHConfigFile(configPath, updatedConfig); err != nil {
		return nil, err
	}

	authStatus := model.SSHHostAuthAuthorized
	if !s.testWithKey(req.HostName, req.User, port, privateKey) {
		authStatus = model.SSHHostAuthFailed
	}

	host := model.SSHHost{
		Alias:      req.Alias,
		HostName:   req.HostName,
		User:       req.User,
		Port:       port,
		AuthStatus: authStatus,
	}
	if err := sshHostRepo.Create(&host); err != nil {
		_ = rollbackSSHConfigHost(configPath, configContent, req.Alias)
		return nil, err
	}

	var info dto.SSHHostInfo
	_ = copier.Copy(&info, &host)
	return &info, nil
}

func (s *SSHHostService) Delete(req dto.ForceDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	configPath, err := sshConfigPath()
	if err != nil {
		return err
	}
	configContent, err := readSSHConfigFile(configPath)
	if err != nil {
		return err
	}

	for _, id := range req.IDs {
		host, err := sshHostRepo.Get(repo.WithByID(id))
		if err != nil {
			continue
		}
		configContent, err = ssh.RemoveSSHConfigHost(configContent, host.Alias)
		if err != nil {
			return err
		}
	}
	if err := writeSSHConfigFile(configPath, configContent); err != nil {
		return err
	}
	return sshHostRepo.Delete(repo.WithByIDs(req.IDs))
}

func (s *SSHHostService) Test(req dto.SSHHostTest) bool {
	if req.ID != 0 {
		host, err := sshHostRepo.Get(repo.WithByID(req.ID))
		if err != nil {
			return false
		}
		_, privateKey, err := loadDefaultSSHKeyPair()
		if err != nil {
			return false
		}
		port := host.Port
		if port == 0 {
			port = 22
		}
		return s.testWithKey(host.HostName, host.User, port, privateKey)
	}
	if err := validateSSHHostInput(req.Alias, req.HostName, req.User); err != nil {
		return false
	}
	_, privateKey, err := loadDefaultSSHKeyPair()
	if err != nil {
		return false
	}
	port := req.Port
	if port == 0 {
		port = 22
	}
	return s.testWithKey(req.HostName, req.User, port, privateKey)
}

func (s *SSHHostService) testWithKey(hostName, user string, port int, privateKey []byte) bool {
	connInfo := ssh.ConnInfo{
		Addr:       hostName,
		Port:       port,
		User:       user,
		AuthMode:   "key",
		PrivateKey: privateKey,
		DialTimeOut: 10 * time.Second,
	}
	client, err := ssh.NewClient(connInfo)
	if err != nil {
		return false
	}
	client.Close()
	return true
}

func validateSSHHostInput(alias, hostName, user string) error {
	if !sshHostAliasPattern.MatchString(alias) {
		return buserr.New("ErrSSHHostAliasInvalid")
	}
	if !isValidSSHHostName(hostName) {
		return buserr.New("ErrSSHHostNameInvalid")
	}
	if !sshHostUserPattern.MatchString(user) {
		return buserr.New("ErrSSHHostUserInvalid")
	}
	return nil
}

func isValidSSHHostName(hostName string) bool {
	if len(hostName) == 0 || len(hostName) > 255 {
		return false
	}
	if net.ParseIP(hostName) != nil {
		return true
	}
	if strings.HasPrefix(hostName, "[") && strings.HasSuffix(hostName, "]") {
		return net.ParseIP(hostName[1:len(hostName)-1]) != nil
	}
	hostPattern := regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*$`)
	return hostPattern.MatchString(hostName)
}

func decodeSSHHostPassword(encoded string) (string, error) {
	if len(encoded) == 0 {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func sshConfigPath() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(current.HomeDir, ".ssh", "config"), nil
}

func defaultSSHKeyPaths() (privatePath, publicPath string, err error) {
	current, err := user.Current()
	if err != nil {
		return "", "", err
	}
	base := filepath.Join(current.HomeDir, ".ssh", defaultSSHKeyName)
	return base, base + ".pub", nil
}

func readSSHConfigFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

func writeSSHConfigFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), constant.FilePerm)
}

func ensureDefaultSSHKey() error {
	privatePath, publicPath, err := defaultSSHKeyPaths()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(privatePath), 0o700); err != nil {
		return err
	}
	if _, err := os.Stat(privatePath); err == nil {
		if _, pubErr := os.Stat(publicPath); pubErr == nil {
			return nil
		}
	}
	tmpPrivate := privatePath + ".tmp"
	tmpPublic := privatePath + ".tmp.pub"
	cmdMgr := cmd.NewCommandMgr(cmd.WithTimeout(2 * time.Minute))
	if err := cmdMgr.Run("ssh-keygen", "-t", "ed25519", "-f", tmpPrivate, "-N", "", "-q"); err != nil {
		_ = os.Remove(tmpPrivate)
		_ = os.Remove(tmpPublic)
		return fmt.Errorf("generate ssh key failed: %w", err)
	}
	if err := os.Rename(tmpPrivate, privatePath); err != nil {
		return err
	}
	return os.Rename(tmpPublic, publicPath)
}

func loadDefaultSSHKeyPair() (publicKey, privateKey []byte, err error) {
	_, publicPath, err := defaultSSHKeyPaths()
	if err != nil {
		return nil, nil, err
	}
	privatePath, _, err := defaultSSHKeyPaths()
	if err != nil {
		return nil, nil, err
	}
	publicKey, err = os.ReadFile(publicPath)
	if err != nil {
		return nil, nil, err
	}
	privateKey, err = os.ReadFile(privatePath)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

func rollbackSSHConfigHost(configPath, originalContent, alias string) error {
	content, err := ssh.RemoveSSHConfigHost(originalContent, alias)
	if err != nil {
		return err
	}
	return writeSSHConfigFile(configPath, content)
}
