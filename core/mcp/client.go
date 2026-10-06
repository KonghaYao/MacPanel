package mcp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	panelauth "github.com/1Panel-dev/1Panel/core/app/auth"
	"github.com/1Panel-dev/1Panel/core/constant"
	"github.com/1Panel-dev/1Panel/core/global"
)

const (
	defaultRequestTimeout = 10 * time.Minute
	maxPanelResponseBytes = 32 << 20 // 32 MiB
)

// PanelHTTPResponse is the loopback Panel API response.
type PanelHTTPResponse struct {
	StatusCode int
	Body       []byte
}

// PanelClient performs signed loopback requests to the local Panel API.
type PanelClient struct {
	httpClient *http.Client
}

// NewPanelClient creates a loopback client for internal Panel API calls.
func NewPanelClient() *PanelClient {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // loopback to local panel cert
	return &PanelClient{
		httpClient: &http.Client{
			Timeout:   defaultRequestTimeout,
			Transport: transport,
		},
	}
}

// Do executes a signed Panel API request using the MCP authenticated API key.
func (c *PanelClient) Do(ctx context.Context, key *AuthenticatedAPIKey, method, apiPath string, query map[string]string, body []byte) (*PanelHTTPResponse, error) {
	if key == nil {
		return nil, fmt.Errorf("missing authenticated API key")
	}
	if err := validateHTTPMethod(method); err != nil {
		return nil, err
	}
	normalizedPath := normalizeAPIPath(apiPath)
	if err := validateAPIPath(normalizedPath); err != nil {
		return nil, err
	}

	target, err := url.Parse(panelBaseURL() + normalizedPath)
	if err != nil {
		return nil, err
	}
	if len(query) > 0 {
		values := target.Query()
		for name, value := range query {
			values.Set(name, value)
		}
		target.RawQuery = values.Encode()
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(strings.TrimSpace(method)), target.String(), bodyReader)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	applySignedHeaders(req.Header, key)
	if key.ClientIP != "" {
		req.Header.Set("X-Forwarded-For", key.ClientIP)
		req.Header.Set("X-Real-IP", key.ClientIP)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxPanelResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(respBody) > maxPanelResponseBytes {
		return nil, fmt.Errorf("panel API response exceeds %d bytes", maxPanelResponseBytes)
	}
	return &PanelHTTPResponse{StatusCode: resp.StatusCode, Body: respBody}, nil
}

func panelBaseURL() string {
	scheme := "http"
	if global.CONF.Conn.SSL == constant.StatusEnable || global.CONF.Conn.SSL == constant.StatusMux {
		scheme = "https"
	}
	return scheme + "://127.0.0.1:" + global.CONF.Conn.Port
}

func applySignedHeaders(headers http.Header, key *AuthenticatedAPIKey) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	token := panelauth.GenerateMD5("1panel" + key.Secret + timestamp)
	headers.Set("1Panel-Token", token)
	headers.Set("1Panel-Timestamp", timestamp)
	if key.KeyID != "" && key.KeyID != "legacy" {
		headers.Set("1Panel-Key-ID", key.KeyID)
	}
}
