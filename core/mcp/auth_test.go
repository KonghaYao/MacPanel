package mcp

import (
	"context"
	"errors"
	"strings"
	"net/http"
	"testing"

	mcpauth "github.com/modelcontextprotocol/go-sdk/auth"
)

func TestVerifyBearerAPIKey_MissingToken(t *testing.T) {
	_, err := VerifyBearerAPIKey(context.Background(), "", &http.Request{})
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if !errors.Is(err, mcpauth.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
	if !strings.Contains(err.Error(), "missing bearer token") {
		t.Fatalf("unexpected message: %v", err)
	}
}

func TestAPIKeyFromTokenInfo(t *testing.T) {
	key := AuthenticatedAPIKey{KeyID: "k1", KeyName: "test", Kind: "apiKey", Secret: "s"}
	info := &mcpauth.TokenInfo{
		Extra: map[string]any{authenticatedAPIKeyExtraKey: key},
	}
	got, ok := APIKeyFromTokenInfo(info)
	if !ok || got.KeyID != "k1" || got.Secret != "s" {
		t.Fatalf("APIKeyFromTokenInfo: ok=%v got=%+v", ok, got)
	}
	if _, ok := APIKeyFromTokenInfo(nil); ok {
		t.Fatal("nil info should not yield key")
	}
	if _, ok := APIKeyFromTokenInfo(&mcpauth.TokenInfo{}); ok {
		t.Fatal("empty extra should not yield key")
	}
}

