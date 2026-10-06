package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	agentServer "github.com/1Panel-dev/1Panel/agent/server"
	adminCmd "github.com/1Panel-dev/1Panel/core/cmd/server/cmd"
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
	writeDaemonPID()
	defer removeDaemonPID()

	sock := paths.SocketPath()

	go func() {
		agentServer.Start()
	}()

	if err := waitForSocket(sock, 60*time.Second); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v (starting core anyway)\n", err)
	}

	if isDaemonChild() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-sigCh
			removeDaemonPID()
			os.Exit(0)
		}()
	}

	coreServer.Start()
}

func main() {
	var daemon bool

	root := &cobra.Command{
		Use:   "macpanel",
		Short: "MacPanel unified server (core + agent)",
		Run: func(cmd *cobra.Command, args []string) {
			if daemon && !isDaemonChild() {
				if err := startDaemon(); err != nil {
					fmt.Fprintf(os.Stderr, "daemon: %v\n", err)
					os.Exit(1)
				}
				return
			}
			runUnified()
		},
	}
	root.Flags().BoolVarP(&daemon, "daemon", "d", false, "Run in background (alias for 'macpanel start', macOS only)")

	root.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start macpanel in background (macOS only)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := startDaemon(); err != nil {
				fmt.Fprintf(os.Stderr, "start: %v\n", err)
				os.Exit(1)
			}
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stop background macpanel (macOS only)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := stopDaemon(); err != nil {
				fmt.Fprintf(os.Stderr, "stop: %v\n", err)
				os.Exit(1)
			}
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "restart",
		Short: "Restart background macpanel (macOS only)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := restartDaemon(); err != nil {
				fmt.Fprintf(os.Stderr, "restart: %v\n", err)
				os.Exit(1)
			}
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show background macpanel status (macOS only)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := printDaemonStatus(); err != nil {
				fmt.Fprintf(os.Stderr, "status: %v\n", err)
				os.Exit(1)
			}
		},
	})

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

	adminCmd.MountAdminCommands(root)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
