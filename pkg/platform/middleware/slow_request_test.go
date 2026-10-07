package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestShouldSkipSlowRequest_StaticPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want bool
	}{
		{"/favicon.ico", true},
		{"/", true},
		{"/assets/app.js", true},
		{"/public/logo.png", true},
		{"/api/v2/static/config.json", true},
		{"/api/v2/images/theme.png", true},
		{"/api/v2/health/check", true},
		{"/1panel/swagger/index.html", true},
		{"/api/v2/containers/search", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodGet, tt.path, nil)
			if got := shouldSkipSlowRequest(c); got != tt.want {
				t.Fatalf("shouldSkipSlowRequest(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestShouldSkipSlowRequest_SSE(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v2/containers/search/log?follow=true", nil)
	c.Request.Header.Set("Accept", "text/event-stream")

	if !shouldSkipSlowRequest(c) {
		t.Fatal("expected SSE log stream request to be skipped")
	}
}

func TestShouldSkipSlowRequest_FollowLogWithoutSSEAccept(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v2/containers/search/log?follow=true", nil)

	if !shouldSkipSlowRequest(c) {
		t.Fatal("expected follow=true log stream path to be skipped")
	}
}

func TestShouldSkipSlowRequest_FollowOnNonLogPath(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v2/containers/search?follow=true", nil)

	if shouldSkipSlowRequest(c) {
		t.Fatal("follow=true on non-log path should not be skipped")
	}
}

func TestIsEventStreamRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		accept string
		want   bool
	}{
		{"text/event-stream", true},
		{"application/json, text/event-stream", true},
		{"TEXT/EVENT-STREAM", true},
		{"application/json", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.accept, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			if got := isEventStreamRequest(req); got != tt.want {
				t.Fatalf("isEventStreamRequest(%q) = %v, want %v", tt.accept, got, tt.want)
			}
		})
	}
}

func TestSlowRequest_SkipRequestCallback(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	logged := false
	router.Use(SlowRequest(SlowRequestConfig{
		Component: "test",
		Threshold: 1,
		Logf: func(format string, args ...any) {
			logged = true
		},
		SkipRequest: func(c *gin.Context) bool {
			return c.Request.URL.Path == "/api/v2/containers/list"
		},
	}))
	router.GET("/api/v2/containers/list", func(c *gin.Context) {
		time.Sleep(5 * time.Millisecond)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v2/containers/list", nil)
	router.ServeHTTP(httptest.NewRecorder(), req)

	if logged {
		t.Fatal("SkipRequest callback should prevent slow request logging")
	}
}
