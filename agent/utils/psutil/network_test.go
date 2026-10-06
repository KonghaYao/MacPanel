//go:build darwin

package psutil

import (
	"context"
	"testing"
)

func TestConnectionsWithContextDarwin(t *testing.T) {
	conns, err := ConnectionsWithContext(context.Background(), "all")
	if err != nil {
		t.Fatalf("ConnectionsWithContext returned error: %v", err)
	}
	if len(conns) == 0 {
		t.Fatal("expected at least one connection on darwin")
	}
}
