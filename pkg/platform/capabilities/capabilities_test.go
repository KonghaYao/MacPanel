package capabilities

import (
	"runtime"
	"testing"
)

func TestDarwinDisablesSSHConfig(t *testing.T) {
	if !IsDarwin() {
		t.Skip("darwin-only test")
	}
	features := darwinFeatures()
	if features.SshdConfig {
		t.Fatal("sshd_config must be disabled on darwin")
	}
}

func TestLinuxEnablesSSHConfig(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("linux-only test")
	}
	features := linuxFeatures()
	if !features.SshdConfig {
		t.Fatal("sshd_config must be enabled on linux")
	}
}

func TestCurrentReturnsDisabledSSHOnDarwin(t *testing.T) {
	if !IsDarwin() {
		t.Skip("darwin-only test")
	}
	caps := Current()
	if caps.Features.SshdConfig {
		t.Fatal("platform capabilities must expose sshd_config=false on darwin")
	}
}
