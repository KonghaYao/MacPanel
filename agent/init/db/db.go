package db

import (
	"os"
	"path"
	"time"

	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/common"
)

func Init() {
	global.DB = common.LoadDBConnByPath(path.Join(global.Dir.DbDir, "agent.db"), "agent")
	global.TaskDB = common.LoadDBConnByPath(path.Join(global.Dir.DbDir, "task.db"), "task")
	global.MonitorDB = common.LoadDBConnByPath(path.Join(global.Dir.DbDir, "monitor.db"), "monitor")
	global.GPUMonitorDB = common.LoadDBConnByPath(path.Join(global.Dir.DbDir, "gpu_monitor.db"), "gpu_monitor")
	global.VLLMMonitorDB = common.LoadDBConnByPath(path.Join(global.Dir.DbDir, "vllm_monitor.db"), "vllm_monitor")
	global.AlertDB = common.LoadDBConnByPath(path.Join(global.Dir.DbDir, "alert.db"), "alert")

	initCoreDBIfExists()
}

func initCoreDBIfExists() {
	if _, err := os.Stat(path.Join(global.Dir.DbDir, "core.db")); err == nil {
		global.CoreDB = common.LoadDBConnByPath(path.Join(global.Dir.DbDir, "core.db"), "core")
	}
}

// WaitAndInitCoreDB polls until core.db exists (created by core on first boot) and loads it.
func WaitAndInitCoreDB(timeout time.Duration) {
	if global.CoreDB != nil {
		return
	}
	deadline := time.Now().Add(timeout)
	coreDBPath := path.Join(global.Dir.DbDir, "core.db")
	for time.Now().Before(deadline) {
		if _, err := os.Stat(coreDBPath); err == nil {
			global.CoreDB = common.LoadDBConnByPath(coreDBPath, "core")
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}
