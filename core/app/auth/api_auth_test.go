package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetAPIClientIP_LoopbackUsesForwardedClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodGet, "/api/v2/test", nil)
	req.RemoteAddr = "127.0.0.1:54321"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	c.Request = req

	got := GetAPIClientIP(c, "")
	if got != "203.0.113.10" {
		t.Fatalf("GetAPIClientIP() = %q, want 203.0.113.10", got)
	}
}

func TestGetAPIClientIP_NonLoopbackIgnoresForwardedWithoutTrustedProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodGet, "/api/v2/test", nil)
	req.RemoteAddr = "198.51.100.9:54321"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	c.Request = req

	got := GetAPIClientIP(c, "")
	if got != "198.51.100.9" {
		t.Fatalf("GetAPIClientIP() = %q, want 198.51.100.9", got)
	}
}
