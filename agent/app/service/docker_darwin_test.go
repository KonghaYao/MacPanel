//go:build darwin

package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/constant"
)

func TestUpdateConfMirrorsOnDarwinSkipsRestart(t *testing.T) {
	dir := t.TempDir()
	daemonPath := filepath.Join(dir, "daemon.json")
	if err := os.WriteFile(daemonPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("seed daemon.json: %v", err)
	}
	originalPath := constant.DaemonJsonPath
	constant.DaemonJsonPath = daemonPath
	t.Cleanup(func() {
		constant.DaemonJsonPath = originalPath
	})

	err := (&DockerService{}).UpdateConf(dto.SettingUpdate{
		Key:   "Mirrors",
		Value: "https://mirror.example.com",
	}, true)
	if err != nil {
		t.Fatalf("UpdateConf failed: %v", err)
	}

	data, err := os.ReadFile(daemonPath)
	if err != nil {
		t.Fatalf("read daemon.json: %v", err)
	}
	var conf map[string]interface{}
	if err := json.Unmarshal(data, &conf); err != nil {
		t.Fatalf("parse daemon.json: %v", err)
	}
	mirrors, ok := conf["registry-mirrors"].([]interface{})
	if !ok || len(mirrors) != 1 || mirrors[0] != "https://mirror.example.com" {
		t.Fatalf("registry-mirrors = %#v", conf["registry-mirrors"])
	}
}

func TestUpdateConfMirrorsOnDarwinOrbStackPath(t *testing.T) {
	dir := t.TempDir()
	orbConfigDir := filepath.Join(dir, ".orbstack", "config")
	if err := os.MkdirAll(orbConfigDir, 0o755); err != nil {
		t.Fatalf("mkdir orbstack config: %v", err)
	}
	daemonPath := filepath.Join(orbConfigDir, "docker.json")
	if err := os.WriteFile(daemonPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("seed orbstack docker.json: %v", err)
	}
	originalPath := constant.DaemonJsonPath
	constant.DaemonJsonPath = daemonPath
	t.Cleanup(func() {
		constant.DaemonJsonPath = originalPath
	})

	err := (&DockerService{}).UpdateConf(dto.SettingUpdate{
		Key:   "Mirrors",
		Value: "https://mirror.orbstack.example.com",
	}, true)
	if err != nil {
		t.Fatalf("UpdateConf failed: %v", err)
	}

	data, err := os.ReadFile(daemonPath)
	if err != nil {
		t.Fatalf("read orbstack docker.json: %v", err)
	}
	var conf map[string]interface{}
	if err := json.Unmarshal(data, &conf); err != nil {
		t.Fatalf("parse orbstack docker.json: %v", err)
	}
	mirrors, ok := conf["registry-mirrors"].([]interface{})
	if !ok || len(mirrors) != 1 || mirrors[0] != "https://mirror.orbstack.example.com" {
		t.Fatalf("registry-mirrors = %#v", conf["registry-mirrors"])
	}
}
