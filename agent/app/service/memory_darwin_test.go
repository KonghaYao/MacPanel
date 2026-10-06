//go:build darwin

package service

import (
	"os"
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

func TestDarwinProcessHelpers(t *testing.T) {
	pid := int32(os.Getpid())
	mem, err := getMemoryDetail(pid)
	if err != nil {
		t.Fatalf("getMemoryDetail: %v", err)
	}
	if mem.RSS == 0 {
		t.Fatal("expected non-zero RSS for current process")
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		t.Fatalf("NewProcess: %v", err)
	}
	envs, err := getProcessEnviron(p)
	if err != nil {
		t.Fatalf("getProcessEnviron: %v", err)
	}
	if len(envs) == 0 {
		t.Fatal("expected environment variables")
	}
}
