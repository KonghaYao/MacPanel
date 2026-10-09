package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/creack/pty"
)

// Close must only signal the shell process; anything written to the pty is
// forwarded by a foreground program (a tmux client, for example) to the
// processes it manages and would kill the user's session on disconnect.
func TestLocalCommandCloseDoesNotWriteToPty(t *testing.T) {
	received := filepath.Join(t.TempDir(), "received.bin")
	// -isig -icanon keeps the line discipline from turning \x03 into SIGINT and
	// \x04 into EOF, so anything written to the pty is recorded verbatim. Real
	// full screen programs such as tmux put the terminal into raw mode too.
	cmd := exec.Command("/bin/sh", "-c", `stty -isig -icanon && exec cat > "$RECEIVED"`)
	cmd.Env = append(os.Environ(), "RECEIVED="+received)
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("start command: %v", err)
	}
	waitForCondition(t, 5*time.Second, func() bool {
		_, statErr := os.Stat(received)
		return statErr == nil
	})

	lcmd := &LocalCommand{cmd: cmd, pty: ptyFile}
	if err := lcmd.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	waitProcessExit(t, cmd, 5*time.Second)

	data, err := os.ReadFile(received)
	if err != nil {
		t.Fatalf("read received: %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("Close wrote %q to the pty", data)
	}
}

// Closing a terminal session must leave a detached tmux session running: the
// shell we started dies, the tmux server (own session) keeps its panes alive.
func TestLocalCommandCloseKeepsTmuxSession(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping tmux integration test in short mode")
	}
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		t.Skip("tmux is not installed")
	}

	// tmux puts its socket at $TMUX_TMPDIR/tmux-<uid>/<name> and a unix socket
	// path is limited to ~104 bytes, so t.TempDir() is too long here.
	tmpDir := filepath.Join("/tmp", fmt.Sprintf("macpanel-terminal-test-%d", os.Getpid()))
	if err := os.MkdirAll(tmpDir, 0o700); err != nil {
		t.Fatalf("create tmux tmpdir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	socket := fmt.Sprintf("localcmd-test-%d", os.Getpid())
	session := func() bool {
		check := exec.Command(tmuxPath, "-L", socket, "has-session", "-t", "t")
		check.Env = tmuxEnv(tmpDir)
		return check.Run() == nil
	}

	cmd := exec.Command("/bin/sh")
	cmd.Env = tmuxEnv(tmpDir)
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("start shell: %v", err)
	}
	lcmd := &LocalCommand{cmd: cmd, pty: ptyFile}
	if _, err := fmt.Fprintf(ptyFile, "%s -L %s new-session -s t\n", tmuxPath, socket); err != nil {
		t.Fatalf("write tmux command: %v", err)
	}
	waitForCondition(t, 15*time.Second, session)

	if err := lcmd.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	waitProcessExit(t, cmd, 5*time.Second)
	time.Sleep(time.Second)

	if !session() {
		kill := exec.Command(tmuxPath, "-L", socket, "kill-server")
		kill.Env = tmuxEnv(tmpDir)
		_ = kill.Run()
		t.Fatalf("Close killed the tmux session")
	}
	kill := exec.Command(tmuxPath, "-L", socket, "kill-server")
	kill.Env = tmuxEnv(tmpDir)
	_ = kill.Run()
}

func tmuxEnv(tmpDir string) []string {
	return append(os.Environ(), "TERM=xterm-256color", "TMUX_TMPDIR="+tmpDir)
}

func waitForCondition(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("condition not met within %s", timeout)
}

func waitProcessExit(t *testing.T, cmd *exec.Cmd, timeout time.Duration) {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatalf("process did not exit within %s", timeout)
	}
}
