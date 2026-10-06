package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/model"
)

type rustFSConnConfig struct {
	AppInstallID uint
	Name         string
	Status       string
	Version      string
	Endpoint     string
	AccessKey    string
	SecretKey    string
	APIPort      int
	ConsolePort  int
	UseSSL       bool
	PathStyle    bool
}

func parseRustFSInstall(install model.AppInstall) (rustFSConnConfig, error) {
	cfg := rustFSConnConfig{
		AppInstallID: install.ID,
		Name:         install.Name,
		Status:       install.Status,
		Version:      install.Version,
		PathStyle:    true,
		UseSSL:       false,
		ConsolePort:  install.HttpPort,
	}
	if strings.TrimSpace(install.Env) == "" {
		return cfg, fmt.Errorf("rustfs install %s has empty env", install.Name)
	}
	envMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(install.Env), &envMap); err != nil {
		return cfg, fmt.Errorf("parse rustfs env of %s: %w", install.Name, err)
	}
	access := envValue(envMap, "RUSTFS_ACCESS_KEY")
	secret := envValue(envMap, "RUSTFS_SECRET_KEY")
	apiPortStr := envValue(envMap, "PANEL_APP_PORT_API")
	if access == "" || secret == "" || apiPortStr == "" {
		return cfg, fmt.Errorf("rustfs install %s is missing RUSTFS_ACCESS_KEY, RUSTFS_SECRET_KEY or PANEL_APP_PORT_API", install.Name)
	}
	apiPort, err := parsePort(apiPortStr)
	if err != nil {
		return cfg, fmt.Errorf("rustfs install %s has invalid PANEL_APP_PORT_API %q", install.Name, apiPortStr)
	}
	cfg.AccessKey = access
	cfg.SecretKey = secret
	cfg.APIPort = apiPort
	cfg.Endpoint = fmt.Sprintf("127.0.0.1:%d", apiPort)
	if httpPort := envValue(envMap, "PANEL_APP_PORT_HTTP"); httpPort != "" {
		if p, err := parsePort(httpPort); err == nil {
			cfg.ConsolePort = p
		}
	}
	return cfg, nil
}

func parsePort(raw string) (int, error) {
	port, err := strconvAtoiStrict(raw)
	if err != nil {
		return 0, err
	}
	if port <= 0 || port > 65535 {
		return 0, fmt.Errorf("port out of range")
	}
	return port, nil
}

func strconvAtoiStrict(raw string) (int, error) {
	var n int
	if _, err := fmt.Sscanf(strings.TrimSpace(raw), "%d", &n); err != nil {
		return 0, err
	}
	if fmt.Sprintf("%d", n) != strings.TrimSpace(raw) {
		return 0, fmt.Errorf("not an integer")
	}
	return n, nil
}

func readLimited(r io.Reader, limit int64) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, r, limit+1); err != nil && err != io.EOF {
		return nil, err
	}
	if int64(buf.Len()) > limit {
		return nil, fmt.Errorf("object exceeds preview size limit")
	}
	return buf.Bytes(), nil
}
