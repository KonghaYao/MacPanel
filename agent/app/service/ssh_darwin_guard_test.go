//go:build darwin

package service

import (
	"testing"

	"github.com/1Panel-dev/1Panel/pkg/platform/capabilities"
)

func TestOperateSSHRejectedOnDarwin(t *testing.T) {
	err := NewISSHService().OperateSSH("enable")
	if err == nil {
		t.Fatal("OperateSSH(enable) must fail on darwin")
	}
	if !capabilities.IsNotSupportedOnMac(err) {
		t.Fatalf("expected NOT_SUPPORTED_ON_DARWIN, got %v", err)
	}
}

func TestGetSSHInfoRejectedOnDarwin(t *testing.T) {
	_, err := NewISSHService().GetSSHInfo()
	if err == nil {
		t.Fatal("GetSSHInfo must fail on darwin")
	}
	if !capabilities.IsNotSupportedOnMac(err) {
		t.Fatalf("expected NOT_SUPPORTED_ON_DARWIN, got %v", err)
	}
}
