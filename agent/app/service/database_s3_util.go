package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const (
	s3SourceRustFS  = "rustfs"
	s3SourceManual  = "manual"
	s3PreviewImage  = "image"
	s3PreviewText   = "text"
	s3PreviewNone   = "none"
	s3PreviewImgMax = 8 << 20
	s3PreviewTxtMax = 512 << 10
)

func envValue(env map[string]interface{}, key string) string {
	v, ok := env[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case json.Number:
		return strings.TrimSpace(t.String())
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func splitS3Endpoint(raw string, useSSL bool) (host string, ssl bool, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false, fmt.Errorf("endpoint is required")
	}
	ssl = useSSL
	if strings.Contains(raw, "://") {
		u, parseErr := url.Parse(raw)
		if parseErr != nil {
			return "", false, fmt.Errorf("invalid endpoint: %w", parseErr)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return "", false, fmt.Errorf("endpoint scheme must be http or https")
		}
		if u.Host == "" {
			return "", false, fmt.Errorf("endpoint host is required")
		}
		return u.Host, u.Scheme == "https", nil
	}
	return raw, ssl, nil
}

func normalizeS3Prefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	prefix = strings.ReplaceAll(prefix, "\\", "/")
	prefix = strings.TrimPrefix(prefix, "/")
	if prefix == "" || prefix == "/" {
		return ""
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return prefix
}

func joinS3Key(prefix, name string) string {
	p := strings.Trim(strings.ReplaceAll(prefix, "\\", "/"), "/")
	n := strings.Trim(strings.ReplaceAll(name, "\\", "/"), "/")
	if p == "" {
		return n
	}
	if n == "" {
		return p
	}
	return p + "/" + n
}

func s3ObjectName(key, prefix string) string {
	current := normalizeS3Prefix(prefix)
	name := strings.TrimPrefix(key, current)
	return strings.TrimSuffix(name, "/")
}

func isS3FolderKey(key string) bool {
	return strings.HasSuffix(key, "/")
}

func validateS3FolderName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("folder name is required")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("folder name must not contain '/'")
	}
	return nil
}

func s3PreviewKind(contentType, key string) string {
	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	ext := strings.ToLower(path.Ext(key))
	switch {
	case strings.HasPrefix(ct, "image/"), isImageExt(ext):
		return s3PreviewImage
	case strings.HasPrefix(ct, "text/"), isTextContentType(ct), isTextExt(ext):
		return s3PreviewText
	default:
		return s3PreviewNone
	}
}

func isImageExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".svg", ".ico":
		return true
	default:
		return false
	}
}

func isTextContentType(ct string) bool {
	switch ct {
	case "application/json", "application/xml", "application/javascript", "application/x-yaml", "application/yaml":
		return true
	default:
		return false
	}
}

func isTextExt(ext string) bool {
	switch ext {
	case ".txt", ".md", ".json", ".xml", ".yml", ".yaml", ".csv", ".log", ".ini", ".conf", ".js", ".ts", ".css", ".html", ".htm", ".go", ".py", ".sh", ".toml":
		return true
	default:
		return false
	}
}

func s3PreviewLimit(kind string) int64 {
	switch kind {
	case s3PreviewImage:
		return s3PreviewImgMax
	case s3PreviewText:
		return s3PreviewTxtMax
	default:
		return 0
	}
}
