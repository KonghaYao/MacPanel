package psutil

import (
	"context"

	"github.com/1Panel-dev/1Panel/pkg/platform/capabilities"
	"github.com/shirou/gopsutil/v4/net"
)

const defaultConnectionLimit = 32768

// ConnectionsWithContext returns network connections for the requested kind.
// Darwin implements ConnectionsMax as not implemented in gopsutil; use Connections there.
func ConnectionsWithContext(ctx context.Context, kind string) ([]net.ConnectionStat, error) {
	if capabilities.IsDarwin() {
		return net.ConnectionsWithContext(ctx, kind)
	}
	return net.ConnectionsMaxWithContext(ctx, kind, defaultConnectionLimit)
}
