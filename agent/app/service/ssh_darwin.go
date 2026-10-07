//go:build darwin

package service

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/pkg/platform/capabilities"
)

const darwinSSHDLaunchdLabel = "com.openssh.sshd"

var darwinSSHEnableOperations = map[string]struct{}{
	"start":   {},
	"restart": {},
	"enable":  {},
}

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

func rejectDarwinSSHEnableOperation(operation string) error {
	if _, blocked := darwinSSHEnableOperations[operation]; blocked {
		return capabilities.ErrNotSupportedOnMac
	}
	return nil
}

func skipDarwinSSHServiceRestart() bool {
	return true
}

func loadDarwinSSHServiceStatus(data *dto.SSHInfo) {
	data.IsExist = true
	data.IsActive = isDarwinRemoteLoginRunning()
	data.AutoStart = isDarwinRemoteLoginEnabled()
	if !data.IsActive {
		data.Message = "Remote Login is disabled on macOS"
	}
}

func operateDarwinSSH(operation string) error {
	switch operation {
	case "stop", "disable":
		if !isDarwinRemoteLoginRunning() && !isDarwinRemoteLoginEnabled() {
			return nil
		}
		if err := disableDarwinRemoteLogin(); err != nil {
			return fmt.Errorf("disable remote login failed, err: %v", err)
		}
		return nil
	default:
		return capabilities.ErrNotSupportedOnMac
	}
}

func isDarwinRemoteLoginRunning() bool {
	output, err := exec.Command("launchctl", "print", "system/"+darwinSSHDLaunchdLabel).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "state = running")
}

func isDarwinRemoteLoginEnabled() bool {
	output, err := exec.Command("systemsetup", "-getremotelogin").CombinedOutput()
	if err != nil {
		return isDarwinRemoteLoginRunning()
	}
	return strings.Contains(strings.ToLower(string(output)), "on")
}

func disableDarwinRemoteLogin() error {
	if output, err := exec.Command("systemsetup", "-setremotelogin", "off").CombinedOutput(); err == nil {
		_ = output
		return nil
	}
	return exec.Command("launchctl", "bootout", "system/system/"+darwinSSHDLaunchdLabel).Run()
}
