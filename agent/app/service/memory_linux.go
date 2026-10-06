//go:build linux

package service

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func getMemoryDetail(pid int32) (*MemoryDetail, error) {
	mem := &MemoryDetail{}

	if err := readStatus(pid, mem); err != nil {
		return nil, err
	}

	if err := readSmapsRollup(pid, mem); err != nil {
		if err := readSmaps(pid, mem); err != nil {
			return nil, err
		}
	}
	return mem, nil
}

func getProcessEnviron(p *process.Process) ([]string, error) {
	return p.Environ()
}

func getProcessOpenFiles(p *process.Process) ([]process.OpenFilesStat, error) {
	return p.OpenFiles()
}

func readStatus(pid int32, mem *MemoryDetail) error {
	filePath := fmt.Sprintf("/proc/%d/status", pid)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024

		switch key {
		case "VmRSS":
			mem.RSS = value
		case "VmSize":
			mem.VMS = value
		case "VmData":
			mem.Data = value
		case "VmSwap":
			mem.Swap = value
		case "VmExe":
			mem.Text = value
		case "RssShmem":
			mem.Shared = value
		case "VmHWM":
			mem.HWM = value
		case "VmStk":
			mem.Stack = value
		case "VmLck":
			mem.Locked = value
		}
	}

	return scanner.Err()
}

func readSmapsRollup(pid int32, mem *MemoryDetail) error {
	filePath := fmt.Sprintf("/proc/%d/smaps_rollup", pid)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024

		switch key {
		case "Pss":
			mem.PSS = value
		case "Private_Clean", "Private_Dirty":
			mem.USS += value
		case "Shared_Clean", "Shared_Dirty":
			if mem.Shared == 0 {
				mem.Shared = value
			}
		}
	}

	return scanner.Err()
}

func readSmaps(pid int32, mem *MemoryDetail) error {
	filePath := fmt.Sprintf("/proc/%d/smaps", pid)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024

		switch key {
		case "Pss":
			mem.PSS += value
		case "Private_Clean", "Private_Dirty":
			mem.USS += value
		case "Shared_Clean", "Shared_Dirty":
			if mem.Shared == 0 {
				mem.Shared += value
			}
		}
	}

	return scanner.Err()
}
