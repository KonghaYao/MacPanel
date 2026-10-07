package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// SlowRequestLogPrefix is a fixed marker for grepping slow HTTP logs.
const SlowRequestLogPrefix = "[MacPanel-SLOW-HTTP]"

const defaultSlowRequestThreshold = 500 * time.Millisecond

// SlowRequestConfig configures slow-request logging middleware.
type SlowRequestConfig struct {
	Component string
	Threshold time.Duration
	Logf      func(format string, args ...any)
	// SkipRequest optionally skips logging for requests handled elsewhere (e.g. core proxy to agent).
	SkipRequest func(*gin.Context) bool
}

// SlowRequest logs requests whose total handler time exceeds the threshold.
// Threshold defaults to 500ms and can be overridden with MACPANEL_SLOW_REQUEST_MS (0 disables).
func SlowRequest(cfg SlowRequestConfig) gin.HandlerFunc {
	threshold := cfg.Threshold
	if threshold == 0 {
		threshold = slowRequestThresholdFromEnv()
	}

	component := cfg.Component
	if component == "" {
		component = "unknown"
	}

	return func(c *gin.Context) {
		if threshold <= 0 || shouldSkipSlowRequest(c) || websocket.IsWebSocketUpgrade(c.Request) || (cfg.SkipRequest != nil && cfg.SkipRequest(c)) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		latency := time.Since(start)
		if latency < threshold || cfg.Logf == nil {
			return
		}

		query := c.Request.URL.RawQuery
		if len(query) > 200 {
			query = query[:200] + "..."
		}

		cfg.Logf(
			"%s component=%s method=%s path=%s query=%q status=%d latency_ms=%d client_ip=%s",
			SlowRequestLogPrefix,
			component,
			c.Request.Method,
			c.Request.URL.Path,
			query,
			c.Writer.Status(),
			latency.Milliseconds(),
			c.ClientIP(),
		)
	}
}

func slowRequestThresholdFromEnv() time.Duration {
	raw := strings.TrimSpace(os.Getenv("MACPANEL_SLOW_REQUEST_MS"))
	if raw == "" {
		return defaultSlowRequestThreshold
	}

	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 0 {
		return defaultSlowRequestThreshold
	}
	if ms == 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

func shouldSkipSlowRequest(c *gin.Context) bool {
	path := c.Request.URL.Path
	switch {
	case path == "/favicon.ico", path == "/":
		return true
	case strings.HasPrefix(path, "/assets/"):
		return true
	case strings.HasPrefix(path, "/public/"):
		return true
	case strings.HasPrefix(path, "/api/v2/static/"):
		return true
	case strings.HasPrefix(path, "/api/v2/images/"):
		return true
	case path == "/api/v2/health/check":
		return true
	case strings.HasPrefix(path, "/1panel/swagger"):
		return true
	}

	if isEventStreamRequest(c.Request) {
		return true
	}

	if c.Query("follow") == "true" && strings.HasSuffix(path, "/search/log") {
		return true
	}

	return false
}

func isEventStreamRequest(req *http.Request) bool {
	return strings.Contains(strings.ToLower(req.Header.Get("Accept")), "text/event-stream")
}
