//go:build !darwin

package service

import (
	"fmt"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto/response"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
)

func (s *DiskService) GetCompleteDiskInfo() (*response.CompleteDiskInfo, error) {
	var diskInfos []response.DiskBasicInfo
	cmdMgr := cmd.NewCommandMgr(cmd.WithTimeout(20 * time.Second))
	output, err := cmdMgr.RunWithStdout("lsblk", "-J", "-o", "NAME,SIZE,TYPE,MOUNTPOINT,FSTYPE,MODEL,SERIAL,TRAN,ROTA")
	if err == nil {
		diskInfos, err = parseLsblkJsonOutput(output)
		if err == nil {
			result := organizeDiskInfo(diskInfos)
			return &result, nil
		}
	}
	output, err = cmdMgr.RunWithStdout("lsblk", "-P", "-o", "NAME,SIZE,TYPE,MOUNTPOINT,FSTYPE,MODEL,SERIAL,TRAN,ROTA")
	if err != nil {
		return nil, fmt.Errorf("failed to run lsblk command: %v", err)
	}
	diskInfos, err = parseLsblkOutput(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lsblk output: %v", err)
	}
	result := organizeDiskInfo(diskInfos)
	return &result, nil
}
