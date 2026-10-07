//go:build darwin

package service

import (
	"testing"

	"github.com/1Panel-dev/1Panel/pkg/platform/capabilities"
)

func TestOperateSSHEnableRejectedOnDarwin(t *testing.T) {
	for _, operation := range []string{"enable", "start", "restart"} {
		err := NewISSHService().OperateSSH(operation)
		if err == nil {
			t.Fatalf("OperateSSH(%s) must fail on darwin", operation)
		}
		if !capabilities.IsNotSupportedOnMac(err) {
			t.Fatalf("OperateSSH(%s): expected NOT_SUPPORTED_ON_DARWIN, got %v", operation, err)
		}
	}
}

func TestOperateSSHDisableAllowedOnDarwin(t *testing.T) {
	if err := NewISSHService().OperateSSH("disable"); err != nil {
		t.Fatalf("OperateSSH(disable) should succeed on darwin, got %v", err)
	}
}

func TestGetSSHInfoAllowedOnDarwin(t *testing.T) {
	info, err := NewISSHService().GetSSHInfo()
	if err != nil {
		t.Fatalf("GetSSHInfo should succeed on darwin, got %v", err)
	}
	if info == nil {
		t.Fatal("GetSSHInfo returned nil")
	}
	if !info.IsExist {
		t.Fatal("GetSSHInfo should report SSH as available on darwin")
	}
}
