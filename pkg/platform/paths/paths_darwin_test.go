//go:build darwin

package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBootstrapCreatesDefaultConfig(t *testing.T) {
	ResetBootstrapState()
	home := t.TempDir()
	t.Setenv("MACPANEL_HOME", home)

	if err := Bootstrap("v9.9.9-test"); err != nil {
		t.Fatalf("Bootstrap: %v", err)
	}

	configPath := filepath.Join(home, "config", "1pctl")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	content := string(data)
	for _, want := range []string{
		"BASE_DIR=",
		"ORIGINAL_PORT=9999",
		"ORIGINAL_VERSION=v9.9.9-test",
		"ORIGINAL_USERNAME=admin",
		"ORIGINAL_PASSWORD=",
		"LANGUAGE=zh",
	} {
		if !containsLine(content, want) {
			t.Fatalf("config missing %q:\n%s", want, content)
		}
	}

	// Second call must be a no-op and keep the same file.
	if err := Bootstrap("v0.0.0-should-not-apply"); err != nil {
		t.Fatalf("Bootstrap second call: %v", err)
	}
	data2, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config again: %v", err)
	}
	if string(data2) != content {
		t.Fatalf("Bootstrap rewrote existing config")
	}
}

func containsLine(content, prefix string) bool {
	for _, line := range splitLines(content) {
		if len(line) >= len(prefix) && line[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func splitLines(content string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			line := content[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(content) {
		lines = append(lines, content[start:])
	}
	return lines
}
