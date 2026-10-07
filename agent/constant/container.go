package constant

import "github.com/1Panel-dev/1Panel/pkg/platform/paths"

const (
	ContainerOpStart   = "start"
	ContainerOpStop    = "stop"
	ContainerOpRestart = "restart"
	ContainerOpKill    = "kill"
	ContainerOpPause   = "pause"
	ContainerOpUnpause = "unpause"
	ContainerOpRename  = "rename"
	ContainerOpRemove  = "remove"

	ComposeOpStop    = "stop"
	ComposeOpRestart = "restart"
	ComposeOpRemove  = "remove"
)

var DaemonJsonPath = "/etc/docker/daemon.json"

func InitDaemonJsonPath(dockerBinaryPath string) {
	DaemonJsonPath = paths.ResolveDockerDaemonJsonPath(dockerBinaryPath)
}
