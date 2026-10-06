package mcp

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

const apiV2Prefix = "/api/v2"

var blockedExactPaths = map[string]struct{}{
	"/api/v2/hosts/terminal/local":     {},
	"/api/v2/hosts/terminal/ssh":       {},
	"/api/v2/hosts/terminal/container": {},
}

var blockedPathPrefixes = []string{
	"/api/v2/internal/",
}

var allowedInternalPaths = map[string]struct{}{}

// BlockedPaths returns the configured blocked API paths and prefixes.
func BlockedPaths() (exact []string, prefixes []string) {
	exact = make([]string, 0, len(blockedExactPaths))
	for p := range blockedExactPaths {
		exact = append(exact, p)
	}
	return exact, append([]string(nil), blockedPathPrefixes...)
}

func isBlockedAPIPath(apiPath string) bool {
	normalized, err := canonicalAPIPath(apiPath)
	if err != nil {
		return false
	}
	return isBlockedNormalizedPath(normalized)
}

func isBlockedNormalizedPath(normalized string) bool {
	if _, blocked := blockedExactPaths[normalized]; blocked {
		return true
	}
	for _, prefix := range blockedPathPrefixes {
		if strings.HasPrefix(normalized, prefix) {
			if _, allowed := allowedInternalPaths[normalized]; allowed {
				continue
			}
			return true
		}
	}
	return false
}

func normalizeAPIPath(apiPath string) string {
	normalized, err := canonicalAPIPath(apiPath)
	if err != nil {
		return ""
	}
	return normalized
}

func canonicalAPIPath(apiPath string) (string, error) {
	raw := strings.TrimSpace(apiPath)
	if raw == "" {
		return "", errInvalidAPIPath
	}
	if strings.ContainsAny(raw, "?#") {
		return "", errInvalidAPIPath
	}
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return "", errInvalidAPIPath
	}
	if strings.Contains(decoded, "\x00") {
		return "", errInvalidAPIPath
	}
	cleaned := path.Clean("/" + decoded)
	if cleaned == "/" {
		return "", errInvalidAPIPath
	}
	if !strings.HasPrefix(cleaned, apiV2Prefix) {
		if strings.HasPrefix(cleaned, "/") {
			cleaned = apiV2Prefix + cleaned
		} else {
			cleaned = apiV2Prefix + "/" + cleaned
		}
	}
	return cleaned, nil
}

func validateAPIPath(apiPath string) error {
	normalized, err := canonicalAPIPath(apiPath)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(normalized, apiV2Prefix+"/") {
		return errInvalidAPIPath
	}
	if isBlockedNormalizedPath(normalized) {
		return errBlockedAPIPath
	}
	return nil
}

func validateHTTPMethod(method string) error {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead:
		return nil
	default:
		return errInvalidHTTPMethod
	}
}
