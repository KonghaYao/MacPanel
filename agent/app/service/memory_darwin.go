//go:build darwin

package service

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
	"golang.org/x/sys/unix"
)

func getMemoryDetail(pid int32) (*MemoryDetail, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		return nil, err
	}

	return &MemoryDetail{
		RSS: memInfo.RSS,
		VMS: memInfo.VMS,
	}, nil
}

func getProcessEnviron(p *process.Process) ([]string, error) {
	return darwinProcessEnviron(p.Pid)
}

func getProcessOpenFiles(p *process.Process) ([]process.OpenFilesStat, error) {
	return darwinProcessOpenFiles(p.Pid)
}

func darwinProcessEnviron(pid int32) ([]string, error) {
	procargs, nargs, err := darwinProcArgs(pid)
	if err != nil {
		return nil, err
	}

	chunks := bytes.Split(procargs, []byte{0})
	if len(chunks) <= 1 {
		return nil, nil
	}

	i := 1
	for ; i < len(chunks) && len(chunks[i]) == 0; i++ {
	}
	if nargs > len(chunks)-i {
		nargs = len(chunks) - i
	}
	i += nargs

	envs := make([]string, 0, len(chunks)-i)
	for ; i < len(chunks); i++ {
		if len(chunks[i]) == 0 {
			continue
		}
		envs = append(envs, string(chunks[i]))
	}
	return envs, nil
}

func darwinProcArgs(pid int32) ([]byte, int, error) {
	procargs, err := unix.SysctlRaw("kern.procargs2", int(pid))
	if err != nil {
		return nil, 0, err
	}
	if len(procargs) < 4 {
		return nil, 0, fmt.Errorf("invalid procargs2 buffer for pid %d", pid)
	}
	nargs := int(binary.LittleEndian.Uint32(procargs[:4]))
	return procargs[4:], nargs, nil
}

func darwinProcessOpenFiles(pid int32) ([]process.OpenFilesStat, error) {
	output, err := exec.Command("lsof", "-n", "-P", "-p", strconv.Itoa(int(pid)), "-F", "pn").Output()
	if err != nil {
		return nil, err
	}

	var files []process.OpenFilesStat
	var current process.OpenFilesStat
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		switch line[0] {
		case 'p':
			if current.Path != "" {
				files = append(files, current)
			}
			current = process.OpenFilesStat{}
		case 'n':
			current.Path = line[1:]
		case 'f':
			if fd, err := strconv.ParseUint(line[1:], 10, 64); err == nil {
				current.Fd = fd
			}
		}
	}
	if current.Path != "" {
		files = append(files, current)
	}
	return files, nil
}
