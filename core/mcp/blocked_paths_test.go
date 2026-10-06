package mcp

import (
	"errors"
	"net/http"
	"testing"
)

func TestNormalizeAPIPath(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"/hosts", "/api/v2/hosts"},
		{"hosts", "/api/v2/hosts"},
		{"/api/v2/hosts", "/api/v2/hosts"},
		{"  /api/v2/settings  ", "/api/v2/settings"},
	}
	for _, tt := range tests {
		if got := normalizeAPIPath(tt.in); got != tt.want {
			t.Errorf("normalizeAPIPath(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIsBlockedAPIPath(t *testing.T) {
	blocked := []string{
		"/api/v2/hosts/terminal/local",
		"/api/v2/hosts/terminal/ssh",
		"/api/v2/internal/anything",
	}
	for _, p := range blocked {
		if !isBlockedAPIPath(p) {
			t.Errorf("expected blocked: %q", p)
		}
	}
	allowed := []string{
		"/api/v2/hosts",
		"/api/v2/settings/search",
	}
	for _, p := range allowed {
		if isBlockedAPIPath(p) {
			t.Errorf("expected allowed: %q", p)
		}
	}
}

func TestValidateAPIPath(t *testing.T) {
	if err := validateAPIPath("/api/v2/settings/search"); err != nil {
		t.Fatalf("allowed path: %v", err)
	}
	if err := validateAPIPath("/api/v2/hosts/terminal/local"); !errors.Is(err, errBlockedAPIPath) {
		t.Fatalf("blocked path err = %v, want errBlockedAPIPath", err)
	}
	if err := validateAPIPath("/"); !errors.Is(err, errInvalidAPIPath) {
		t.Fatalf("root path err = %v, want errInvalidAPIPath", err)
	}
}

func TestValidateAPIPathRejectsQueryAndEncodedBypass(t *testing.T) {
	bypassAttempts := []string{
		"/api/v2/hosts/terminal/local?terminalRevalidate=1",
		"/api/v2/internal%2Fterminal/sessions",
		"/api/v2/internal%2fsettings",
	}
	for _, p := range bypassAttempts {
		err := validateAPIPath(p)
		if errors.Is(err, errBlockedAPIPath) || errors.Is(err, errInvalidAPIPath) {
			continue
		}
		t.Fatalf("path %q: err = %v, want blocked or invalid", p, err)
	}
}

func TestBlockedPathsExport(t *testing.T) {
	exact, prefixes := BlockedPaths()
	if len(exact) != len(blockedExactPaths) {
		t.Fatalf("exact count = %d, want %d", len(exact), len(blockedExactPaths))
	}
	if len(prefixes) != len(blockedPathPrefixes) {
		t.Fatalf("prefix count = %d, want %d", len(prefixes), len(blockedPathPrefixes))
	}
}

func TestValidateHTTPMethod(t *testing.T) {
	for _, m := range []string{"GET", "post", " Put ", "PATCH", "DELETE", "HEAD"} {
		if err := validateHTTPMethod(m); err != nil {
			t.Errorf("method %q: %v", m, err)
		}
	}
	if err := validateHTTPMethod("OPTIONS"); !errors.Is(err, errInvalidHTTPMethod) {
		t.Fatalf("OPTIONS err = %v, want errInvalidHTTPMethod", err)
	}
	if err := validateHTTPMethod(""); !errors.Is(err, errInvalidHTTPMethod) {
		t.Fatalf("empty method err = %v, want errInvalidHTTPMethod", err)
	}
	_ = http.MethodGet
}
