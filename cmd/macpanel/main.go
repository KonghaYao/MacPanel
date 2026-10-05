package main

import (
	"fmt"
	"net"
	"os"
	"time"

	agentServer "github.com/1Panel-dev/1Panel/agent/server"
	coreServer "github.com/1Panel-dev/1Panel/core/server"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
	"github.com/spf13/cobra"
)

func waitForSocket(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("unix", path, time.Second)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("agent socket %q not ready within %s", path, timeout)
}

func runUnified() {
	sock := paths.SocketPath()

	go func() {
		agentServer.Start()
	}()

	if err := waitForSocket(sock, 60*time.Second); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v (starting core anyway)\n", err)
	}

	coreServer.Start()
}

func main() {
	root := &cobra.Command{
		Use:   "macpanel",
		Short: "MacPanel unified server (core + agent)",
		Run: func(cmd *cobra.Command, args []string) {
			runUnified()
		},
	}

	root.AddCommand(&cobra.Command{
		Use:   "core",
		Short: "Start core server only",
		Run: func(cmd *cobra.Command, args []string) {
			coreServer.Start()
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "agent",
		Short: "Start agent server only",
		Run: func(cmd *cobra.Command, args []string) {
			agentServer.Start()
		},
	})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
