package mcp

import (
	"strings"
	"testing"
)

func TestSearchOpenAPIPaths(t *testing.T) {
	matches, err := searchOpenAPIPaths("accounts")
	if err != nil {
		t.Fatalf("searchOpenAPIPaths: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected matches for keyword accounts")
	}
	for _, m := range matches {
		if m.Path == "" || m.Method == "" {
			t.Fatalf("invalid match: %+v", m)
		}
		text := strings.ToLower(m.Path + " " + m.Method + " " + m.Summary)
		if !strings.Contains(text, "accounts") {
			t.Fatalf("match does not contain keyword: %+v", m)
		}
	}
}

func TestSearchOpenAPIPathsEmptyKeyword(t *testing.T) {
	_, err := searchOpenAPIPaths("   ")
	if err == nil || !strings.Contains(err.Error(), "keyword is required") {
		t.Fatalf("empty keyword err = %v", err)
	}
}

func TestGetOpenAPISchema(t *testing.T) {
	schema, err := getOpenAPISchema("/api/v2/ai/accounts", "POST")
	if err != nil {
		t.Fatalf("getOpenAPISchema: %v", err)
	}
	if schema["path"] != "/api/v2/ai/accounts" {
		t.Fatalf("path = %v", schema["path"])
	}
	if schema["method"] != "POST" {
		t.Fatalf("method = %v", schema["method"])
	}
	op, ok := schema["operation"].(map[string]interface{})
	if !ok || len(op) == 0 {
		t.Fatalf("missing operation schema: %+v", schema)
	}
}

func TestGetOpenAPISchemaNotFound(t *testing.T) {
	_, err := getOpenAPISchema("/api/v2/does-not-exist-mcp-test", "GET")
	if err == nil || !strings.Contains(err.Error(), "path not found") {
		t.Fatalf("err = %v", err)
	}
}
