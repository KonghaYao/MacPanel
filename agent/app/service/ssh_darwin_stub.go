//go:build !darwin

package service

import "github.com/1Panel-dev/1Panel/agent/app/dto"

func EnsureSSHDisabledOnDarwin() {}

func rejectDarwinSSHEnableOperation(operation string) error { return nil }

func skipDarwinSSHServiceRestart() bool { return false }

func loadDarwinSSHServiceStatus(data *dto.SSHInfo) {}

func operateDarwinSSH(operation string) error { return nil }
