//go:build darwin

package viper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
)

func TestInitUsesStableModeWithoutDevConfig(t *testing.T) {
	paths.ResetBootstrapState()
	home := t.TempDir()
	t.Setenv("MACPANEL_HOME", home)

	if err := paths.Bootstrap("v0.3.2-test"); err != nil {
		t.Fatalf("Bootstrap: %v", err)
	}

	Init()

	if got := global.CONF.Base.Mode; got != "stable" {
		t.Fatalf("Mode = %q, want stable", got)
	}
	if got := global.CONF.LogConfig.Level; got != "info" {
		t.Fatalf("Log level = %q, want info", got)
	}
}

func TestInitKeepsDevModeWithDevConfig(t *testing.T) {
	paths.ResetBootstrapState()
	home := t.TempDir()
	t.Setenv("MACPANEL_HOME", home)

	if err := paths.Bootstrap("v0.3.2-test"); err != nil {
		t.Fatalf("Bootstrap: %v", err)
	}

	devConf := filepath.Join(home, "1panel", "conf", "app.yaml")
	if err := os.MkdirAll(filepath.Dir(devConf), 0o700); err != nil {
		t.Fatalf("mkdir conf: %v", err)
	}
	content := []byte(`base:
  mode: dev
  install_dir: /tmp/macpanel-dev
conn:
  port: "8888"
log:
  level: debug
`)
	if err := os.WriteFile(devConf, content, 0o600); err != nil {
		t.Fatalf("write dev conf: %v", err)
	}

	Init()

	if got := global.CONF.Base.Mode; got != "dev" {
		t.Fatalf("Mode = %q, want dev", got)
	}
	if got := global.CONF.Conn.Port; got != "8888" {
		t.Fatalf("Port = %q, want 8888", got)
	}
}
