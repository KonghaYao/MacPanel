package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

var defaultPanelClient = NewPanelClient()

type panelRequestInput struct {
	Method string            `json:"method" jsonschema:"HTTP method such as GET or POST"`
	Path   string            `json:"path" jsonschema:"Panel API path starting with /api/v2/"`
	Body   string            `json:"body,omitempty" jsonschema:"Optional JSON request body as a string"`
	Query  map[string]string `json:"query,omitempty" jsonschema:"Optional query parameters"`
}

type panelRequestOutput struct {
	StatusCode int             `json:"statusCode"`
	Body       json.RawMessage `json:"body"`
}

type panelOpenAPISearchInput struct {
	Keyword string `json:"keyword" jsonschema:"Keyword to search in OpenAPI paths and summaries"`
}

type panelOpenAPISearchOutput struct {
	Matches []openAPIPathMatch `json:"matches"`
}

type panelOpenAPISchemaInput struct {
	Path   string `json:"path" jsonschema:"Panel API path, with or without /api/v2 prefix"`
	Method string `json:"method" jsonschema:"HTTP method such as GET or POST"`
}

type panelOpenAPISchemaOutput struct {
	Path      string                 `json:"path"`
	Method    string                 `json:"method"`
	Operation map[string]interface{} `json:"operation"`
}

// RegisterTools registers generic Panel MCP tools on the server.
func RegisterTools(server *mcpsdk.Server) {
	if server == nil {
		return
	}
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "panel_request",
		Description: "Call a MacPanel API endpoint under /api/v2 using the authenticated API key.",
	}, handlePanelRequest)
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "panel_openapi_search",
		Description: "Search embedded OpenAPI paths by keyword.",
	}, handlePanelOpenAPISearch)
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "panel_openapi_schema",
		Description: "Return OpenAPI operation schema for a Panel API path and method.",
	}, handlePanelOpenAPISchema)
}

func handlePanelRequest(ctx context.Context, req *mcpsdk.CallToolRequest, input panelRequestInput) (*mcpsdk.CallToolResult, panelRequestOutput, error) {
	key, err := authenticatedKeyFromCallTool(req)
	if err != nil {
		return toolError(err), panelRequestOutput{}, nil
	}
	var body []byte
	if strings.TrimSpace(input.Body) != "" {
		if !json.Valid([]byte(input.Body)) {
			return toolError(fmt.Errorf("body must be valid JSON")), panelRequestOutput{}, nil
		}
		body = []byte(input.Body)
	}
	resp, err := defaultPanelClient.Do(ctx, key, input.Method, input.Path, input.Query, body)
	if err != nil {
		return toolError(err), panelRequestOutput{}, nil
	}
	return nil, panelRequestOutput{
		StatusCode: resp.StatusCode,
		Body:       normalizeResponseBody(resp.Body),
	}, nil
}

func handlePanelOpenAPISearch(_ context.Context, _ *mcpsdk.CallToolRequest, input panelOpenAPISearchInput) (*mcpsdk.CallToolResult, panelOpenAPISearchOutput, error) {
	matches, err := searchOpenAPIPaths(input.Keyword)
	if err != nil {
		return toolError(err), panelOpenAPISearchOutput{}, nil
	}
	return nil, panelOpenAPISearchOutput{Matches: matches}, nil
}

func handlePanelOpenAPISchema(_ context.Context, _ *mcpsdk.CallToolRequest, input panelOpenAPISchemaInput) (*mcpsdk.CallToolResult, panelOpenAPISchemaOutput, error) {
	schema, err := getOpenAPISchema(input.Path, input.Method)
	if err != nil {
		return toolError(err), panelOpenAPISchemaOutput{}, nil
	}
	operation, _ := schema["operation"].(map[string]interface{})
	return nil, panelOpenAPISchemaOutput{
		Path:      fmt.Sprint(schema["path"]),
		Method:    fmt.Sprint(schema["method"]),
		Operation: operation,
	}, nil
}

func normalizeResponseBody(body []byte) json.RawMessage {
	if len(body) == 0 {
		return json.RawMessage("null")
	}
	if json.Valid(body) {
		return json.RawMessage(body)
	}
	wrapped, err := json.Marshal(map[string]string{"raw": string(body)})
	if err != nil {
		return json.RawMessage("null")
	}
	return json.RawMessage(wrapped)
}

func authenticatedKeyFromCallTool(req *mcpsdk.CallToolRequest) (*AuthenticatedAPIKey, error) {
	if req == nil || req.Extra == nil || req.Extra.TokenInfo == nil {
		return nil, fmt.Errorf("missing authenticated API key")
	}
	key, ok := APIKeyFromTokenInfo(req.Extra.TokenInfo)
	if !ok || key == nil {
		return nil, fmt.Errorf("missing authenticated API key")
	}
	return key, nil
}

func toolError(err error) *mcpsdk.CallToolResult {
	return &mcpsdk.CallToolResult{
		IsError: true,
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: err.Error()}},
	}
}
