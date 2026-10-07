package mirrors

import (
	"sort"
	"strings"
)

func updateAssignments(content string, updates map[string]string) string {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	seen := map[string]bool{}
	var out []string
	for _, line := range lines {
		key, _, ok := splitAssign(line)
		value, managed := updates[key]
		if !ok || !managed {
			out = append(out, line)
			continue
		}
		if seen[key] || value == "" {
			continue
		}
		seen[key] = true
		out = append(out, key+"="+value)
	}
	keys := make([]string, 0, len(updates))
	for key := range updates {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if updates[key] == "" || seen[key] {
			continue
		}
		out = append(out, key+"="+updates[key])
	}
	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, "\n") + "\n"
}

func readAssignments(content string, keys []string) map[string]string {
	wanted := map[string]struct{}{}
	for _, key := range keys {
		wanted[key] = struct{}{}
	}
	found := map[string]string{}
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		key, value, ok := splitAssign(line)
		if !ok {
			continue
		}
		if _, yes := wanted[key]; yes {
			if _, exists := found[key]; !exists {
				found[key] = value
			}
		}
	}
	return found
}

func splitAssign(line string) (string, string, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ";") {
		return "", "", false
	}
	eq := strings.Index(trimmed, "=")
	if eq <= 0 {
		return "", "", false
	}
	return strings.TrimSpace(trimmed[:eq]), strings.TrimSpace(trimmed[eq+1:]), true
}

func assignmentFileMeaningful(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) != "" {
			return true
		}
	}
	return false
}
