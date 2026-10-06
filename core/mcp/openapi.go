package mcp

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/1Panel-dev/1Panel/core/cmd/server/docs"
)

type swaggerDocument struct {
	BasePath string                            `json:"basePath"`
	Paths    map[string]map[string]interface{} `json:"paths"`
}

type openAPIPathMatch struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Summary string `json:"summary,omitempty"`
}

var (
	swaggerOnce sync.Once
	swaggerDoc  swaggerDocument
	swaggerErr  error
)

func loadSwagger() (swaggerDocument, error) {
	swaggerOnce.Do(func() {
		raw := docs.SwaggerInfo.ReadDoc()
		swaggerErr = json.Unmarshal([]byte(raw), &swaggerDoc)
	})
	return swaggerDoc, swaggerErr
}

func searchOpenAPIPaths(keyword string) ([]openAPIPathMatch, error) {
	doc, err := loadSwagger()
	if err != nil {
		return nil, err
	}
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil, fmt.Errorf("keyword is required")
	}

	matches := make([]openAPIPathMatch, 0)
	for swaggerPath, methods := range doc.Paths {
		fullPath := fullAPIPath(doc.BasePath, swaggerPath)
		for method, operation := range methods {
			method = strings.ToLower(strings.TrimSpace(method))
			if method == "" || method[0] == '$' {
				continue
			}
			summary := operationSummary(operation)
			searchText := strings.ToLower(strings.Join([]string{fullPath, method, summary}, " "))
			if !strings.Contains(searchText, keyword) {
				continue
			}
			matches = append(matches, openAPIPathMatch{
				Path:    fullPath,
				Method:  strings.ToUpper(method),
				Summary: summary,
			})
		}
	}
	return matches, nil
}

func getOpenAPISchema(apiPath, method string) (map[string]interface{}, error) {
	doc, err := loadSwagger()
	if err != nil {
		return nil, err
	}
	method = strings.ToLower(strings.TrimSpace(method))
	if method == "" {
		return nil, fmt.Errorf("method is required")
	}

	swaggerPath := stripAPIPrefix(apiPath)
	methods, ok := doc.Paths[swaggerPath]
	if !ok {
		return nil, fmt.Errorf("path not found in openapi docs")
	}
	operation, ok := methods[method]
	if !ok {
		return nil, fmt.Errorf("method not found for path")
	}
	opMap, ok := operation.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid operation schema")
	}
	return map[string]interface{}{
		"path":      fullAPIPath(doc.BasePath, swaggerPath),
		"method":    strings.ToUpper(method),
		"operation": opMap,
	}, nil
}

func fullAPIPath(basePath, swaggerPath string) string {
	basePath = strings.TrimSuffix(strings.TrimSpace(basePath), "/")
	if basePath == "" {
		basePath = apiV2Prefix
	}
	if strings.HasPrefix(swaggerPath, apiV2Prefix) {
		return swaggerPath
	}
	if !strings.HasPrefix(swaggerPath, "/") {
		swaggerPath = "/" + swaggerPath
	}
	return basePath + swaggerPath
}

func stripAPIPrefix(apiPath string) string {
	normalized := normalizeAPIPath(apiPath)
	return strings.TrimPrefix(normalized, apiV2Prefix)
}

func operationSummary(operation interface{}) string {
	opMap, ok := operation.(map[string]interface{})
	if !ok {
		return ""
	}
	summary, _ := opMap["summary"].(string)
	if summary != "" {
		return summary
	}
	description, _ := opMap["description"].(string)
	return description
}
