//go:build darwin

package i18n

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
)

func TestGetMsgByKeyForCmdWithoutExistingConfig(t *testing.T) {
	paths.ResetBootstrapState()
	home := t.TempDir()
	t.Setenv("MACPANEL_HOME", home)

	global.I18nForCmd = nil

	msg := GetMsgByKeyForCmd("UpdateUser")
	if msg == "" || msg == "UpdateUser" {
		t.Fatalf("expected localized UpdateUser, got %q", msg)
	}

	configPath := filepath.Join(home, "config", "1pctl")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected bootstrap to create config: %v", err)
	}
}
