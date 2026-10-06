//go:build darwin

package websocket

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
)

func loadPlatformSSHSessions(ctx context.Context) ([]sshSession, bool) {
	output, err := exec.CommandContext(ctx, "who", "-u").Output()
	if err != nil {
		return nil, false
	}

	var sessions []sshSession
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		session, ok := parseDarwinWhoLine(line)
		if !ok {
			continue
		}
		sessions = append(sessions, session)
	}
	return sessions, true
}

func parseDarwinWhoLine(line string) (sshSession, bool) {
	hostStart := strings.LastIndex(line, "(")
	hostEnd := strings.LastIndex(line, ")")
	if hostStart == -1 || hostEnd <= hostStart {
		return sshSession{}, false
	}

	host := strings.TrimSpace(line[hostStart+1 : hostEnd])
	prefix := strings.TrimSpace(line[:hostStart])
	fields := strings.Fields(prefix)
	if len(fields) < 4 {
		return sshSession{}, false
	}

	pid, err := strconv.ParseInt(fields[len(fields)-1], 10, 32)
	if err != nil {
		return sshSession{}, false
	}

	loginTime := strings.Join(fields[2:len(fields)-2], " ")
	return sshSession{
		Username:  fields[0],
		Terminal:  fields[1],
		Host:      host,
		PID:       int32(pid),
		LoginTime: loginTime,
	}, true
}
