package ssh

import (
	"bufio"
	"fmt"
	"strings"
)

const managedMarker = "# Managed by MacPanel"

type ConfigHostBlock struct {
	Alias    string
	HostName string
	User     string
	Start    int
	End      int
	Managed  bool
}

func ParseSSHConfig(content string) []ConfigHostBlock {
	lines := strings.Split(content, "\n")
	var blocks []ConfigHostBlock
	var current *ConfigHostBlock
	managedNext := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == managedMarker {
			managedNext = true
			continue
		}
		if strings.HasPrefix(strings.ToLower(trimmed), "host ") {
			if current != nil {
				current.End = i
				blocks = append(blocks, *current)
			}
			alias := strings.TrimSpace(trimmed[5:])
			current = &ConfigHostBlock{
				Alias:   alias,
				Start:   i,
				Managed: managedNext,
			}
			managedNext = false
			continue
		}
		if current == nil || len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}
		key := strings.ToLower(parts[0])
		value := strings.TrimSpace(strings.Join(parts[1:], " "))
		switch key {
		case "hostname":
			current.HostName = value
		case "user":
			current.User = value
		}
	}
	if current != nil {
		current.End = len(lines)
		blocks = append(blocks, *current)
	}
	return blocks
}

func HasSSHConfigHost(content, alias string) bool {
	for _, block := range ParseSSHConfig(content) {
		if block.Alias == alias {
			return true
		}
	}
	return false
}

func RemoveSSHConfigHost(content, alias string) (string, error) {
	lines := strings.Split(content, "\n")
	blocks := ParseSSHConfig(content)
	var target *ConfigHostBlock
	for i := range blocks {
		if blocks[i].Alias == alias {
			target = &blocks[i]
			break
		}
	}
	if target == nil {
		return content, nil
	}

	start := target.Start
	if start > 0 && strings.TrimSpace(lines[start-1]) == managedMarker {
		start--
	}

	var kept []string
	for i, line := range lines {
		if i >= start && i < target.End {
			continue
		}
		kept = append(kept, line)
	}
	result := strings.Join(kept, "\n")
	result = strings.TrimRight(result, "\n")
	if len(result) > 0 {
		result += "\n"
	}
	return result, nil
}

func AddSSHConfigHost(content, alias, hostName, user string) (string, error) {
	if HasSSHConfigHost(content, alias) {
		return "", fmt.Errorf("host alias %q already exists in ssh config", alias)
	}

	block := fmt.Sprintf("%s\nHost %s\n  HostName %s\n  User %s\n", managedMarker, alias, hostName, user)
	content = strings.TrimRight(content, "\n")
	if len(strings.TrimSpace(content)) == 0 {
		return block, nil
	}
	return content + "\n\n" + block, nil
}

func FormatSSHConfigBlock(alias, hostName, user string) string {
	var b strings.Builder
	b.WriteString(managedMarker)
	b.WriteByte('\n')
	b.WriteString("Host ")
	b.WriteString(alias)
	b.WriteByte('\n')
	b.WriteString("  HostName ")
	b.WriteString(hostName)
	b.WriteByte('\n')
	b.WriteString("  User ")
	b.WriteString(user)
	b.WriteByte('\n')
	return b.String()
}

func ScanSSHConfigHosts(content string) []ConfigHostBlock {
	return ParseSSHConfig(content)
}

func ReadSSHConfigLines(content string) []string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
