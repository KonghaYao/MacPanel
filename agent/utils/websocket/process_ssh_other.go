//go:build !darwin

package websocket

import "context"

func loadPlatformSSHSessions(_ context.Context) ([]sshSession, bool) {
	return nil, false
}
