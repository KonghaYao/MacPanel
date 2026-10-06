package service

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto/request"
	"github.com/1Panel-dev/1Panel/agent/utils/common"
	agentPsutil "github.com/1Panel-dev/1Panel/agent/utils/psutil"
	"github.com/1Panel-dev/1Panel/agent/utils/websocket"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessService struct{}

type IProcessService interface {
	StopProcess(req request.ProcessReq) error
	GetProcessInfoByPID(pid int32) (*websocket.PsProcessData, error)
	GetListeningProcess(c context.Context) ([]ListeningProcess, error)
}

func NewIProcessService() IProcessService {
	return &ProcessService{}
}

func (ps *ProcessService) StopProcess(req request.ProcessReq) error {
	proc, err := process.NewProcess(req.PID)
	if err != nil {
		return err
	}
	if err := proc.Kill(); err != nil {
		return err
	}
	return nil
}

type ListeningProcess struct {
	PID      int32
	Port     map[uint32]struct{}
	Protocol uint32
	Name     string
}

func (ps *ProcessService) GetListeningProcess(c context.Context) ([]ListeningProcess, error) {
	conn, err := agentPsutil.ConnectionsWithContext(c, "inet")
	if err != nil {
		return nil, err
	}
	// One cache entry per (PID, socket type) so TCP and UDP sockets are not merged under one Protocol.
	type procKey struct {
		pid      int32
		protocol uint32
	}
	procCache := make(map[procKey]ListeningProcess, 64)

	for _, conn := range conn {
		if conn.Pid == 0 {
			continue
		}

		if (conn.Status == "LISTEN" && conn.Type == syscall.SOCK_STREAM) || (conn.Type == syscall.SOCK_DGRAM && conn.Raddr.Port == 0) {
			key := procKey{pid: conn.Pid, protocol: conn.Type}
			if _, exists := procCache[key]; !exists {
				proc, err := process.NewProcess(conn.Pid)
				if err != nil {
					continue
				}
				procData := ListeningProcess{
					PID: conn.Pid,
				}
				procData.Name, _ = proc.Name()
				procData.Port = make(map[uint32]struct{})
				procData.Port[conn.Laddr.Port] = struct{}{}
				procData.Protocol = conn.Type
				procCache[key] = procData
			} else {
				p := procCache[key]
				p.Port[conn.Laddr.Port] = struct{}{}
				procCache[key] = p
			}
		}
	}

	procs := make([]ListeningProcess, 0, len(procCache))
	for _, proc := range procCache {
		procs = append(procs, proc)
	}

	return procs, nil
}

func (ps *ProcessService) GetProcessInfoByPID(pid int32) (*websocket.PsProcessData, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("get process info by pid %v: %v", pid, err)
	}

	exists, err := p.IsRunning()
	if err != nil || !exists {
		return nil, fmt.Errorf("process %v is not running", pid)
	}

	data := &websocket.PsProcessData{
		PID: pid,
	}

	if name, err := p.Name(); err == nil {
		data.Name = name
	}

	if ppid, err := p.Ppid(); err == nil {
		data.PPID = ppid
	}

	if username, err := p.Username(); err == nil {
		data.Username = username
	}

	if status, err := p.Status(); err == nil {
		if len(status) > 0 {
			data.Status = status[0]
		}
	}

	if createTime, err := agentPsutil.NewProcessCreateTimeResolver().CreateTime(p); err == nil {
		data.StartTime = time.Unix(createTime/1000, 0).Format("2006-01-02 15:04:05")
	}

	if numThreads, err := p.NumThreads(); err == nil {
		data.NumThreads = numThreads
	}

	if connections, err := p.Connections(); err == nil {
		data.NumConnections = len(connections)

		var connects []websocket.ProcessConnect
		for _, conn := range connections {
			pc := websocket.ProcessConnect{
				Status: conn.Status,
				Laddr:  conn.Laddr,
				Raddr:  conn.Raddr,
				PID:    pid,
				Name:   data.Name,
			}
			connects = append(connects, pc)
		}
		data.Connects = connects
	}

	if cpuPercent, err := p.CPUPercent(); err == nil {
		data.CpuValue = cpuPercent
		data.CpuPercent = fmt.Sprintf("%.2f%%", cpuPercent)
	}

	if ioCounters, err := p.IOCounters(); err == nil {
		data.DiskRead = common.FormatBytes(ioCounters.ReadBytes)
		data.DiskWrite = common.FormatBytes(ioCounters.WriteBytes)
	}

	if cmdline, err := p.Cmdline(); err == nil {
		data.CmdLine = cmdline
	}

	if memDetail, err := getMemoryDetail(p.Pid); err == nil {
		data.Rss = common.FormatBytes(memDetail.RSS)
		data.VMS = common.FormatBytes(memDetail.VMS)
		data.HWM = common.FormatBytes(memDetail.HWM)
		data.Data = common.FormatBytes(memDetail.Data)
		data.Stack = common.FormatBytes(memDetail.Stack)
		data.Locked = common.FormatBytes(memDetail.Locked)
		data.Swap = common.FormatBytes(memDetail.Swap)
		data.Dirty = common.FormatBytes(memDetail.Dirty)
		data.RssValue = memDetail.RSS
		data.PSS = common.FormatBytes(memDetail.PSS)
		data.USS = common.FormatBytes(memDetail.USS)
		data.Shared = common.FormatBytes(memDetail.Shared)
		data.Text = common.FormatBytes(memDetail.Text)
	}

	if envs, err := getProcessEnviron(p); err == nil {
		data.Envs = envs
	}

	if openFiles, err := getProcessOpenFiles(p); err == nil {
		data.OpenFiles = openFiles
	}

	return data, nil
}

type MemoryDetail struct {
	RSS    uint64
	VMS    uint64
	HWM    uint64
	Data   uint64
	Stack  uint64
	Locked uint64
	Swap   uint64

	PSS    uint64
	USS    uint64
	Shared uint64
	Text   uint64
	Dirty  uint64
}
