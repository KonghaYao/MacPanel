package mcp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMCPUnauthorized401(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
