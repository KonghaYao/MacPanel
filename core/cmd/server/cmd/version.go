package cmd

import (
	"fmt"

	"github.com/1Panel-dev/1Panel/core/i18n"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	RunE: func(cmd *cobra.Command, args []string) error {
		i18n.UseI18nForCmd(language)
		if !isRoot() {
			fmt.Println(adminHelperCmd("version"))
			return nil
		}
		version, mode := VersionFromCtl()
		if db, err := loadDBConn("core.db"); err == nil {
			if dbVersion := getSettingByKey(db, "SystemVersion"); dbVersion != "" {
				version = dbVersion
			}
		}

		fmt.Println(i18n.GetMsgByKeyForCmd("SystemVersion") + version)
		if paths.IsDarwin() {
			fmt.Println(i18n.GetMsgByKeyForCmd("SystemMode") + mode)
		}
		return nil
	},
}
