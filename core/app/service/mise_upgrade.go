package service

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func miseUpgradeArgs() [][]string {
	return [][]string{
		{"mise", "use", "-g", "github:KonghaYao/MacPanel[bin=macpanel]@latest"},
		{"mise", "reshim"},
	}
}

func (u *UpgradeService) UpgradeByMise() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var parts []string
	for _, args := range miseUpgradeArgs() {
		command := exec.CommandContext(ctx, args[0], args[1:]...)
		out, err := command.CombinedOutput()
		if text := strings.TrimSpace(string(out)); text != "" {
			parts = append(parts, text)
		}
		if err != nil {
			joined := strings.Join(parts, "\n")
			if joined == "" {
				return "", err
			}
			return joined, fmt.Errorf("%s: %w", joined, err)
		}
	}
	return strings.Join(parts, "\n"), nil
}
