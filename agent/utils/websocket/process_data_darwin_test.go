//go:build darwin

package websocket

import "testing"

func TestGetProcessDataDarwin(t *testing.T) {
	data, err := getProcessData(PsProcessConfig{})
	if err != nil {
		t.Fatalf("getProcessData: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected process websocket payload")
	}
}

func TestGetNetConnectionsDarwin(t *testing.T) {
	data, err := getNetConnections(NetConfig{})
	if err != nil {
		t.Fatalf("getNetConnections: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected network websocket payload")
	}
}
