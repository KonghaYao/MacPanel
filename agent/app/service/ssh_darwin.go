//go:build darwin

package service

import (
	"os/exec"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/global"
)

const darwinSSHDLaunchdLabel = "com.openssh.sshd"

// EnsureSSHDisabledOnDarwin verifies Remote Login is off on startup.
// MacPanel never enables SSH on macOS; disabling an already-running sshd requires root.
func EnsureSSHDisabledOnDarwin() {
	if !isDarwinRemoteLoginRunning() {
		return
	}
	if err := disableDarwinRemoteLogin(); err != nil {
		global.LOG.Warnf("remote login is enabled on macOS but could not be disabled automatically: %v", err)
		return
	}
	if isDarwinRemoteLoginRunning() {
		global.LOG.Warn("remote login is still enabled on macOS; disable it manually in System Settings")
	}
}

func isDarwinRemoteLoginRunning() bool {
	output, err := exec.Command("launchctl", "print", "system/"+darwinSSHDLaunchdLabel).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "state = running")
}

func disableDarwinRemoteLogin() error {
	if output, err := exec.Command("systemsetup", "-setremotelogin", "off").CombinedOutput(); err == nil {
		_ = output
		return nil
	}
	return exec.Command("launchctl", "bootout", "system/system/"+darwinSSHDLaunchdLabel).Run()
}
