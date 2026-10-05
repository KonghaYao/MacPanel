package cmd

import (
	"os/exec"
	"runtime"
)

func SudoHandleCmd() string {
	if runtime.GOOS == "darwin" {
		return ""
	}
	cmd := exec.Command("sudo", "-n", "ls")
	if err := cmd.Run(); err == nil {
		return "sudo "
	}
	return ""
}

func Which(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
