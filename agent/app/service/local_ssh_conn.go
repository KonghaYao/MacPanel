package service

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
	"runtime"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/copier"
	"github.com/1Panel-dev/1Panel/agent/utils/encrypt"
	"github.com/1Panel-dev/1Panel/agent/utils/ssh"
)

const localSSHKeyName = "id_ed25519_1panel"

// EnsureLocalSSHConn provisions key-based localhost SSH credentials when missing.
// On macOS the panel runs as the current user, so root login must not be assumed.
func EnsureLocalSSHConn() error {
	connItem, _ := settingRepo.GetValueByKey("LocalSSHConn")
	if len(connItem) != 0 {
		if runtime.GOOS != "darwin" || localSSHConnValid(connItem) {
			return nil
		}
		if err := settingRepo.Update("LocalSSHConn", ""); err != nil {
			return err
		}
	}

	currentInfo, err := user.Current()
	if err != nil {
		return err
	}

	keyPath := localSSHKeyPath(currentInfo)
	if err := os.MkdirAll(path.Dir(keyPath), 0o700); err != nil {
		return err
	}
	if _, err := os.Stat(keyPath); err != nil {
		if err := NewISSHService().CreateRootCert(dto.RootCertOperate{
			EncryptionMode: "ed25519",
			Name:           localSSHKeyName,
			Description:    "1Panel Terminal",
		}); err != nil {
			return err
		}
	}

	privateKey, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	connWithKey := ssh.ConnInfo{
		Addr:       "127.0.0.1",
		User:       localSSHUser(currentInfo),
		Port:       22,
		AuthMode:   "key",
		PrivateKey: privateKey,
	}
	client, err := ssh.NewClient(connWithKey)
	if err != nil {
		return err
	}
	client.Close()

	var conn model.LocalConnInfo
	_ = copier.Copy(&conn, &connWithKey)
	conn.PrivateKey = string(privateKey)
	conn.PassPhrase = ""
	localConn, err := json.Marshal(&conn)
	if err != nil {
		return err
	}
	connAfterEncrypt, err := encrypt.StringEncrypt(string(localConn))
	if err != nil {
		return err
	}
	return settingRepo.Update("LocalSSHConn", connAfterEncrypt)
}

func EnsureLocalSSHConnOnDarwin() {
	if runtime.GOOS != "darwin" {
		return
	}
	if err := EnsureLocalSSHConn(); err != nil {
		global.LOG.Warnf("ensure local ssh conn on darwin failed: %v", err)
	}
}

func localSSHKeyPath(currentInfo *user.User) string {
	if currentInfo == nil || len(currentInfo.HomeDir) == 0 {
		return path.Join("/root", ".ssh", localSSHKeyName)
	}
	return path.Join(currentInfo.HomeDir, ".ssh", localSSHKeyName)
}

func localSSHUser(currentInfo *user.User) string {
	if runtime.GOOS == "darwin" {
		if currentInfo != nil && currentInfo.Username != "" {
			return currentInfo.Username
		}
	}
	return "root"
}

func localSSHConnValid(encrypted string) bool {
	connInfoInDB, err := encrypt.StringDecrypt(encrypted)
	if err != nil {
		return false
	}
	var data model.LocalConnInfo
	if err := json.Unmarshal([]byte(connInfoInDB), &data); err != nil {
		return false
	}
	currentInfo, err := user.Current()
	if err != nil {
		return false
	}
	if data.User != localSSHUser(currentInfo) {
		return false
	}
	connInfo := ssh.ConnInfo{
		Addr:       data.Addr,
		Port:       int(data.Port),
		User:       data.User,
		AuthMode:   data.AuthMode,
		Password:   data.Password,
		PrivateKey: []byte(data.PrivateKey),
	}
	if len(data.PassPhrase) != 0 {
		connInfo.PassPhrase = []byte(data.PassPhrase)
	}
	client, err := ssh.NewClient(connInfo)
	if err != nil {
		return false
	}
	client.Close()
	return true
}
