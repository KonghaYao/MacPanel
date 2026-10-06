package mcp

import "errors"

var (
	errInvalidAPIPath   = errors.New("path must start with /api/v2/")
	errBlockedAPIPath   = errors.New("path is blocked for MCP panel_request")
	errInvalidHTTPMethod = errors.New("unsupported HTTP method")
)
