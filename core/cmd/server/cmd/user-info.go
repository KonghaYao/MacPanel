package cmd

import (
	"fmt"
	"os"

	"github.com/1Panel-dev/1Panel/core/i18n"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(userinfoCmd)
}

var userinfoCmd = &cobra.Command{
	Use: "user-info",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isRoot() {
			i18n.UseI18nForCmd(language)
			fmt.Println(adminHelperCmd("user-info"))
			return nil
		}
		return PrintUserInfo(os.Stdout, language)
	},
}
