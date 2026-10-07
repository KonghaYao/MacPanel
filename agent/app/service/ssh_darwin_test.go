//go:build darwin

package service

import "testing"

func TestIsDarwinRemoteLoginRunning(t *testing.T) {
	// Should not panic; result depends on local system state.
	_ = isDarwinRemoteLoginRunning()
}

func TestDisableDarwinRemoteLoginDoesNotPanic(t *testing.T) {
	if !isDarwinRemoteLoginRunning() {
		return
	}
	_ = disableDarwinRemoteLogin()
}
