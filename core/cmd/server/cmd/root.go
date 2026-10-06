package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/1Panel-dev/1Panel/core/i18n"
	"github.com/1Panel-dev/1Panel/core/server"
	"github.com/1Panel-dev/1Panel/core/utils/cmd"
	"github.com/1Panel-dev/1Panel/core/utils/common"
	"github.com/1Panel-dev/1Panel/core/utils/ctl_conf"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var language string

func init() {
	RootCmd.PersistentFlags().StringVarP(&language, "language", "l", "en", "Set the language")
}

var RootCmd = &cobra.Command{
	Use: "1panel",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start()
		return nil
	},
}

type setting struct {
	ID        uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Key       string    `json:"key" gorm:"type:varchar(256);not null;"`
	Value     string    `json:"value" gorm:"type:varchar(256)"`
	About     string    `json:"about" gorm:"type:longText"`
}

func loadDBConn(dbName string) (*gorm.DB, error) {
	baseDir, err := loadBaseDir()
	if err != nil {
		return nil, err
	}

	db, err := common.GetDBWithPath(path.Join(baseDir, "1panel/db", dbName))
	if err != nil {
		return nil, fmt.Errorf("init my db conn failed, err: %v", err)
	}
	return db, nil
}

func loadBaseDir() (string, error) {
	configFile := paths.ConfigFile()
	baseDir, err := ctl_conf.LoadFromFile(configFile, "BASE_DIR")
	if err != nil {
		return "", fmt.Errorf("handle load `BASE_DIR` failed, err: %v", err)
	}
	if len(baseDir) == 0 {
		return "", fmt.Errorf("error `BASE_DIR` find in %s", configFile)
	}
	return baseDir, nil
}

func getSettingByKey(db *gorm.DB, key string) string {
	var setting setting
	_ = db.Where("key = ?", key).First(&setting).Error
	return setting.Value
}

type LoginLog struct{}

func shouldShowInitialPassword(db *gorm.DB) bool {
	logCount := int64(0)
	_ = db.Model(&LoginLog{}).Where("status = ?", "Success").Count(&logCount).Error
	return logCount == 0
}

func setSettingByKey(db *gorm.DB, key, value string) error {
	return db.Model(&setting{}).Where("key = ?", key).Updates(map[string]interface{}{"value": value}).Error
}

var adminCommandNames = []string{
	"user-info", "user-list", "update", "reset", "listen-ip", "version", "restore",
}

func MountAdminCommands(root *cobra.Command) {
	root.PersistentFlags().StringVarP(&language, "language", "l", "en", "Set the language")

	toMove := make([]*cobra.Command, 0, len(adminCommandNames))
	for _, name := range adminCommandNames {
		for _, c := range RootCmd.Commands() {
			if c.Name() == name {
				toMove = append(toMove, c)
				break
			}
		}
	}
	for _, c := range toMove {
		RootCmd.RemoveCommand(c)
		root.AddCommand(c)
	}
}

func cliName() string {
	if paths.IsDarwin() {
		return "macpanel"
	}
	return "1pctl"
}

func adminHelperCmd(subcmd string) string {
	if paths.IsDarwin() {
		return i18n.GetMsgWithMapForCmd("SudoHelper", map[string]interface{}{"cmd": cliName() + " " + subcmd})
	}
	return i18n.GetMsgWithMapForCmd("SudoHelper", map[string]interface{}{"cmd": "sudo " + cliName() + " " + subcmd})
}

func restartPanel() (string, error) {
	mgr := cmd.NewCommandMgr()
	if paths.IsDarwin() {
		binary, err := os.Executable()
		if err != nil {
			binary = paths.CoreBinaryPath()
		}
		return mgr.RunWithStdout(binary, "restart")
	}
	return mgr.RunWithStdout("1pctl", "restart", "core")
}

func isRoot() bool {
	if paths.IsDarwin() {
		return ownsMacPanelDir()
	}
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}

func ownsMacPanelDir() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return false
	}
	for _, dir := range []string{paths.ConfigDir(), paths.DataDir()} {
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return false
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return false
		}
		if stat.Uid != uint32(uid) {
			return false
		}
	}
	return true
}
