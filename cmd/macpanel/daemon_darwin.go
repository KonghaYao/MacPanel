//go:build darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
)

const daemonEnv = "MACPANEL_DAEMON=1"

func isDaemonChild() bool {
	return os.Getenv("MACPANEL_DAEMON") == "1"
}

func startDaemon() error {
	if isDaemonChild() {
		return nil
	}

	if err := os.MkdirAll(paths.RunDir(), 0o700); err != nil {
		return fmt.Errorf("create run dir: %w", err)
	}

	if running, pid, err := existingDaemonPID(); err != nil {
		return err
	} else if running {
		return fmt.Errorf("macpanel already running (pid %d)", pid)
	}

	logPath := paths.DaemonLogFile()
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	args := daemonChildArgs(os.Args[1:])
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), daemonEnv)
	cmd.Stdin = nil
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("start daemon: %w", err)
	}
	_ = logFile.Close()

	pidPath := paths.PidFile()
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(cmd.Process.Pid)), 0o600); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("write pid file: %w", err)
	}

	fmt.Printf("macpanel started in background (pid %d)\n", cmd.Process.Pid)
	fmt.Printf("log: %s\n", logPath)
	fmt.Printf("pid: %s\n", pidPath)
	return nil
}

func daemonChildArgs(args []string) []string {
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg {
		case "-d", "--daemon", "start", "stop", "restart", "status":
			continue
		default:
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

func stopDaemon() error {
	running, pid, err := existingDaemonPID()
	if err != nil {
		return err
	}
	if !running {
		return fmt.Errorf("macpanel is not running")
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		_ = os.Remove(paths.PidFile())
		fmt.Println("macpanel stopped")
		return nil
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		_ = os.Remove(paths.PidFile())
		return fmt.Errorf("send SIGTERM to pid %d: %w", pid, err)
	}

	if !waitForProcessExit(pid, 30*time.Second) {
		_ = proc.Signal(syscall.SIGKILL)
		if !waitForProcessExit(pid, 5*time.Second) {
			return fmt.Errorf("macpanel (pid %d) did not stop within timeout", pid)
		}
	}

	_ = os.Remove(paths.PidFile())
	fmt.Println("macpanel stopped")
	return nil
}

func restartDaemon() error {
	running, _, err := existingDaemonPID()
	if err != nil {
		return err
	}
	if running {
		if err := stopDaemon(); err != nil {
			return err
		}
	}
	return startDaemon()
}

func printDaemonStatus() error {
	running, pid, err := existingDaemonPID()
	if err != nil {
		return err
	}
	if running {
		fmt.Printf("macpanel is running (pid %d)\n", pid)
	} else {
		fmt.Println("macpanel is stopped")
	}
	return nil
}

func waitForProcessExit(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		proc, err := os.FindProcess(pid)
		if err != nil {
			return true
		}
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func existingDaemonPID() (running bool, pid int, err error) {
	pidPath := paths.PidFile()
	data, err := os.ReadFile(pidPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, 0, nil
		}
		return false, 0, fmt.Errorf("read pid file: %w", err)
	}

	pid, err = strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return false, 0, nil
	}
	if pid <= 0 {
		return false, 0, nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false, 0, nil
	}
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		_ = os.Remove(pidPath)
		return false, 0, nil
	}
	return true, pid, nil
}

func writeDaemonPID() {
	if !isDaemonChild() {
		return
	}
	if err := os.MkdirAll(paths.RunDir(), 0o700); err != nil {
		return
	}
	_ = os.WriteFile(paths.PidFile(), []byte(strconv.Itoa(os.Getpid())), 0o600)
}

func removeDaemonPID() {
	if !isDaemonChild() {
		return
	}
	_ = os.Remove(paths.PidFile())
}
