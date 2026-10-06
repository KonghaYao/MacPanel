package service

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/1Panel-dev/1Panel/agent/app/model"
)

func TestEnvValue(t *testing.T) {
	env := map[string]interface{}{
		"RUSTFS_ACCESS_KEY":   " rustfsadmin ",
		"PANEL_APP_PORT_API":  float64(9000),
		"PANEL_APP_PORT_HTTP": json.Number("9001"),
	}
	if got := envValue(env, "RUSTFS_ACCESS_KEY"); got != "rustfsadmin" {
		t.Fatalf("access key = %q", got)
	}
	if got := envValue(env, "PANEL_APP_PORT_API"); got != "9000" {
		t.Fatalf("api port = %q", got)
	}
	if got := envValue(env, "PANEL_APP_PORT_HTTP"); got != "9001" {
		t.Fatalf("http port = %q", got)
	}
	if got := envValue(env, "MISSING"); got != "" {
		t.Fatalf("missing = %q", got)
	}
}

func TestParseRustFSInstall(t *testing.T) {
	env, err := json.Marshal(map[string]interface{}{
		"PANEL_APP_PORT_API":  9000,
		"PANEL_APP_PORT_HTTP": 9001,
		"RUSTFS_ACCESS_KEY":   "ak",
		"RUSTFS_SECRET_KEY":   "sk",
	})
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := parseRustFSInstall(model.AppInstall{
		BaseModel: model.BaseModel{ID: 7},
		Name:      "rustfs",
		Status:    "Running",
		Version:   "1.0.1",
		Env:       string(env),
		HttpPort:  9001,
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Endpoint != "127.0.0.1:9000" {
		t.Fatalf("endpoint = %q", cfg.Endpoint)
	}
	if cfg.APIPort != 9000 || cfg.ConsolePort != 9001 {
		t.Fatalf("ports api=%d console=%d", cfg.APIPort, cfg.ConsolePort)
	}
	if cfg.AccessKey != "ak" || cfg.SecretKey != "sk" || !cfg.PathStyle || cfg.UseSSL {
		t.Fatalf("unexpected cfg %+v", cfg)
	}
}

func TestParseRustFSInstallMissingAPIPort(t *testing.T) {
	env, err := json.Marshal(map[string]interface{}{
		"RUSTFS_ACCESS_KEY": "ak",
		"RUSTFS_SECRET_KEY": "sk",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = parseRustFSInstall(model.AppInstall{Name: "rustfs", Env: string(env)})
	if err == nil || !strings.Contains(err.Error(), "PANEL_APP_PORT_API") {
		t.Fatalf("expected missing API port error, got %v", err)
	}
}

func TestSplitS3Endpoint(t *testing.T) {
	host, ssl, err := splitS3Endpoint("http://127.0.0.1:9000", true)
	if err != nil {
		t.Fatal(err)
	}
	if host != "127.0.0.1:9000" || ssl {
		t.Fatalf("host=%q ssl=%v", host, ssl)
	}
	host, ssl, err = splitS3Endpoint("s3.amazonaws.com", true)
	if err != nil {
		t.Fatal(err)
	}
	if host != "s3.amazonaws.com" || !ssl {
		t.Fatalf("host=%q ssl=%v", host, ssl)
	}
	if _, _, err := splitS3Endpoint("ftp://x", false); err == nil {
		t.Fatal("expected scheme error")
	}
}

func TestS3KeyHelpers(t *testing.T) {
	if got := normalizeS3Prefix("photos"); got != "photos/" {
		t.Fatalf("prefix = %q", got)
	}
	if got := joinS3Key("photos/", "2024"); got != "photos/2024" {
		t.Fatalf("join = %q", got)
	}
	if got := s3ObjectName("photos/2024/", "photos/"); got != "2024" {
		t.Fatalf("name = %q", got)
	}
	if !isS3FolderKey("photos/") {
		t.Fatal("expected folder key")
	}
	if err := validateS3FolderName("a/b"); err == nil {
		t.Fatal("expected nested name error")
	}
}

func TestS3PreviewKind(t *testing.T) {
	if got := s3PreviewKind("image/png", "a.bin"); got != s3PreviewImage {
		t.Fatalf("kind = %s", got)
	}
	if got := s3PreviewKind("", "readme.md"); got != s3PreviewText {
		t.Fatalf("kind = %s", got)
	}
	if got := s3PreviewKind("application/octet-stream", "a.bin"); got != s3PreviewNone {
		t.Fatalf("kind = %s", got)
	}
	if s3PreviewLimit(s3PreviewImage) != s3PreviewImgMax {
		t.Fatal("image limit")
	}
}
