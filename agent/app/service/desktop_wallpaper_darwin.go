//go:build darwin

package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
)

const (
	desktopWallpaperJXA = `ObjC.import("AppKit"); var ws = $.NSWorkspace.sharedWorkspace; var screen = $.NSScreen.mainScreen; var url = ws.desktopImageURLForScreen(screen); url ? url.path.js : ""`
	desktopWallpaperMax = "1920"
)

func (s *DesktopService) GetCompressedWallpaper() ([]byte, error) {
	path, err := macMainDisplayWallpaperPath()
	if err != nil {
		return nil, err
	}
	return compressWallpaperJPEG(path)
}

func macMainDisplayWallpaperPath() (string, error) {
	cmdMgr := cmd.NewCommandMgr()
	output, err := cmdMgr.RunWithStdout("osascript", "-l", "JavaScript", "-e", desktopWallpaperJXA)
	if err != nil {
		return "", fmt.Errorf("read desktop wallpaper path: %w", err)
	}
	path := strings.TrimSpace(output)
	if path == "" {
		return "", fmt.Errorf("desktop wallpaper path is empty")
	}
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("desktop wallpaper path is not absolute")
	}
	clean := filepath.Clean(path)
	if strings.Contains(clean, "..") {
		return "", fmt.Errorf("desktop wallpaper path is invalid")
	}
	info, err := os.Stat(clean)
	if err != nil {
		return "", fmt.Errorf("desktop wallpaper file unavailable: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("desktop wallpaper path is a directory")
	}
	return clean, nil
}

func compressWallpaperJPEG(path string) ([]byte, error) {
	outFile, err := os.CreateTemp("", "macpanel-wallpaper-*.jpg")
	if err != nil {
		return nil, fmt.Errorf("create wallpaper temp file: %w", err)
	}
	outPath := outFile.Name()
	_ = outFile.Close()
	defer os.Remove(outPath)

	cmdMgr := cmd.NewCommandMgr()
	if err := cmdMgr.Run(
		"sips",
		"-s", "format", "jpeg",
		"-s", "formatOptions", "80",
		"-Z", desktopWallpaperMax,
		path,
		"--out", outPath,
	); err != nil {
		return nil, fmt.Errorf("compress desktop wallpaper: %w", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("read compressed desktop wallpaper: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("compressed desktop wallpaper is empty")
	}
	return data, nil
}
