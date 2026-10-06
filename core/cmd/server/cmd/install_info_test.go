//go:build darwin

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
)

func TestPrintInstallInfoFromCtl(t *testing.T) {
	paths.ResetBootstrapState()
	home := t.TempDir()
	t.Setenv("MACPANEL_HOME", home)

	if err := paths.Bootstrap("v0.3.2-test"); err != nil {
		t.Fatalf("Bootstrap: %v", err)
	}

	var buf bytes.Buffer
	PrintInstallInfo(&buf)
	out := buf.String()

	for _, want := range []string{
		"Panel URL: http://127.0.0.1:9999",
		"Username: admin",
		"Password:",
		"Version: v0.3.2-test",
		"macpanel user-info",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

func TestVersionFromCtlStableMode(t *testing.T) {
	paths.ResetBootstrapState()
	home := t.TempDir()
	t.Setenv("MACPANEL_HOME", home)

	if err := paths.Bootstrap("v0.3.2-test"); err != nil {
		t.Fatalf("Bootstrap: %v", err)
	}

	version, mode := VersionFromCtl()
	if version != "v0.3.2-test" {
		t.Fatalf("version = %q, want v0.3.2-test", version)
	}
	if mode != "stable" {
		t.Fatalf("mode = %q, want stable", mode)
	}

	devConf := filepath.Join(home, "1panel", "conf", "app.yaml")
	if err := os.MkdirAll(filepath.Dir(devConf), 0o700); err != nil {
		t.Fatalf("mkdir conf: %v", err)
	}
	if err := os.WriteFile(devConf, []byte("base:\n  mode: dev\n"), 0o600); err != nil {
		t.Fatalf("write dev conf: %v", err)
	}

	_, mode = VersionFromCtl()
	if mode != "dev" {
		t.Fatalf("mode with dev app.yaml = %q, want dev", mode)
	}
}

func TestBuildPanelURL(t *testing.T) {
	if got := buildPanelURL("127.0.0.1", "9999", "entrance"); got != "http://127.0.0.1:9999/entrance" {
		t.Fatalf("unexpected url: %s", got)
	}
	if got := buildPanelURL("127.0.0.1", "9999", ""); got != "http://127.0.0.1:9999" {
		t.Fatalf("unexpected url: %s", got)
	}
}
