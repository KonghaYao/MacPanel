package mcp

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	panelauth "github.com/1Panel-dev/1Panel/core/app/auth"
)

func TestApplySignedHeaders(t *testing.T) {
	key := &AuthenticatedAPIKey{
		KeyID:  "test-key-id",
		Secret: "test-secret",
	}
	headers := make(http.Header)
	before := time.Now().Unix()
	applySignedHeaders(headers, key)
	after := time.Now().Unix()

	ts := headers.Get("1Panel-Timestamp")
	if ts == "" {
		t.Fatal("missing 1Panel-Timestamp")
	}
	parsedTS, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		t.Fatalf("timestamp parse: %v", err)
	}
	if parsedTS < before || parsedTS > after {
		t.Fatalf("timestamp %d out of range [%d,%d]", parsedTS, before, after)
	}

	token := headers.Get("1Panel-Token")
	want := panelauth.GenerateMD5("1panel" + key.Secret + ts)
	if token != want {
		t.Fatalf("token = %q, want %q", token, want)
	}
	if headers.Get("1Panel-Key-ID") != key.KeyID {
		t.Fatalf("Key-ID = %q", headers.Get("1Panel-Key-ID"))
	}
}

func TestApplySignedHeadersLegacyKeyOmitsKeyID(t *testing.T) {
	key := &AuthenticatedAPIKey{KeyID: "legacy", Secret: "legacy-secret"}
	headers := make(http.Header)
	applySignedHeaders(headers, key)
	if headers.Get("1Panel-Key-ID") != "" {
		t.Fatalf("legacy key should not set Key-ID, got %q", headers.Get("1Panel-Key-ID"))
	}
}
