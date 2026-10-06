package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/1Panel-dev/1Panel/core/constant"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/core/i18n"
	"github.com/1Panel-dev/1Panel/core/utils/encrypt"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
	"gorm.io/gorm"
)

func PrintUserInfo(out io.Writer, language string) error {
	i18n.UseI18nForCmd(language)

	db, err := loadDBConn("core.db")
	if err != nil {
		if paths.IsDarwin() {
			return PrintUserInfoFromCtl(out)
		}
		return fmt.Errorf("init my db conn failed, err: %v", err)
	}
	agentDB, err := loadDBConn("agent.db")
	if err != nil {
		return fmt.Errorf("init my agent db conn failed, err: %v", err)
	}

	isEnterpriseVersion := isEnterprise()
	showInitialPassword := shouldShowInitialPassword(db)
	encryptSetting := getSettingByKey(db, "EncryptKey")
	user := ""
	pass := "********"
	if isEnterpriseVersion {
		enterpriseDB, err := loadDBConn("enterprise.db")
		if err != nil {
			return fmt.Errorf("init my enterprise db conn failed, err: %v", err)
		}
		enterpriseUser, enterprisePassword, err := loadEnterpriseSuperAdminInfo(enterpriseDB)
		if err != nil {
			return err
		}
		user = enterpriseUser
		if showInitialPassword {
			pass = enterprisePassword
		}
	} else {
		user = getSettingByKey(db, "UserName")
		if showInitialPassword {
			pass = getSettingByKey(db, "Password")
		}
	}
	if showInitialPassword {
		if len(encryptSetting) == 16 {
			global.CONF.Base.EncryptKey = encryptSetting
			pass, _ = encrypt.StringDecrypt(pass)
		}
	}

	port := getSettingByKey(db, "ServerPort")
	ssl := getSettingByKey(db, "SSL")
	entrance := getSettingByKey(db, "SecurityEntrance")
	address := getSettingByKey(agentDB, "SystemIP")
	domain := getSettingByKey(db, "BindDomain")
	if len(domain) != 0 {
		address = domain
	}

	protocol := "http"
	if ssl == constant.StatusEnable {
		protocol = "https"
	}
	if address == "" {
		address = "$LOCAL_IP"
	}

	fmt.Fprintln(out, i18n.GetMsgByKeyForCmd("UserInfoAddr")+fmt.Sprintf("%s://%s:%s/%s ", protocol, address, port, entrance))
	fmt.Fprintln(out, i18n.GetMsgWithMapForCmd("UpdateUserResult", map[string]interface{}{"name": user}))
	fmt.Fprintln(out, i18n.GetMsgWithMapForCmd("UpdatePasswordResult", map[string]interface{}{"name": pass}))
	updatePasswordCmd := cliName() + " update password"
	if isEnterpriseVersion && strings.TrimSpace(user) != "" {
		updatePasswordCmd += " --username " + user
	}
	fmt.Fprintln(out, i18n.GetMsgByKeyForCmd("UserInfoPassHelp")+updatePasswordCmd)
	fmt.Fprintln(out, i18n.GetMsgWithMapForCmd("UserInfoViewHelp", map[string]interface{}{"cmd": cliName() + " user-info"}))
	return nil
}

func loadEnterpriseSuperAdminInfo(db *gorm.DB) (string, string, error) {
	var user struct {
		Name     string
		Password string
	}
	result := db.Raw("SELECT name, password FROM users WHERE is_super_admin = ? ORDER BY id LIMIT 1", true).Scan(&user)
	if result.Error != nil {
		return "", "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", "", fmt.Errorf("super admin user not found")
	}
	return user.Name, user.Password, nil
}
