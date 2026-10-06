package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Panel-dev/1Panel/core/utils/ctl_conf"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
)

func PrintInstallInfo(out io.Writer) {
	port := ctlConfValue("ORIGINAL_PORT", "9999")
	entrance := strings.Trim(ctl_conf.LoadWithoutPanic("ORIGINAL_ENTRANCE"), `"`)
	username := ctl_conf.LoadWithoutPanic("ORIGINAL_USERNAME")
	password := ctl_conf.LoadWithoutPanic("ORIGINAL_PASSWORD")
	version := ctl_conf.LoadWithoutPanic("ORIGINAL_VERSION")

	panelURL := buildPanelURL("127.0.0.1", port, entrance)
	fmt.Fprintf(out, "Panel URL: %s\n", panelURL)
	if username != "" {
		fmt.Fprintf(out, "Username: %s\n", username)
	}
	if password != "" {
		fmt.Fprintf(out, "Password: %s\n", password)
		fmt.Fprintf(out, "Change password: %s update password\n", cliName())
	}
	if version != "" {
		fmt.Fprintf(out, "Version: %s\n", version)
	}
	fmt.Fprintf(out, "View details: %s user-info\n", cliName())
}

func PrintUserInfoFromCtl(out io.Writer) error {
	port := ctlConfValue("ORIGINAL_PORT", "9999")
	entrance := strings.Trim(ctl_conf.LoadWithoutPanic("ORIGINAL_ENTRANCE"), `"`)
	username := ctl_conf.LoadWithoutPanic("ORIGINAL_USERNAME")
	password := ctl_conf.LoadWithoutPanic("ORIGINAL_PASSWORD")

	panelURL := buildPanelURL("127.0.0.1", port, entrance)
	fmt.Fprintf(out, "Panel URL: %s\n", panelURL)
	if username != "" {
		fmt.Fprintf(out, "Username: %s\n", username)
	}
	if password != "" {
		fmt.Fprintf(out, "Password: %s\n", password)
		fmt.Fprintf(out, "Change password: %s update password\n", cliName())
	}
	return nil
}

func VersionFromCtl() (version, mode string) {
	version = ctl_conf.LoadWithoutPanic("ORIGINAL_VERSION")
	if version == "" {
		version = "unknown"
	}
	if paths.IsDarwin() {
		devConf := filepath.Join(paths.ConfDir(), "app.yaml")
		if _, err := os.Stat(devConf); err == nil {
			return version, "dev"
		}
	}
	return version, "stable"
}

func ctlConfValue(key, fallback string) string {
	value := strings.Trim(ctl_conf.LoadWithoutPanic(key), `"`)
	if value == "" {
		return fallback
	}
	return value
}

func buildPanelURL(host, port, entrance string) string {
	url := fmt.Sprintf("http://%s:%s", host, port)
	entrance = strings.Trim(entrance, "/")
	if entrance != "" {
		url += "/" + entrance
	}
	return url
}
