package mcp

import (
	"net/http"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/gin-gonic/gin"
	mcpauth "github.com/modelcontextprotocol/go-sdk/auth"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

var handler http.Handler

func init() {
	server := newServer()
	streamable := mcpsdk.NewStreamableHTTPHandler(func(*http.Request) *mcpsdk.Server {
		return server
	}, &mcpsdk.StreamableHTTPOptions{
		Stateless: true,
	})
	handler = mcpauth.RequireBearerToken(VerifyBearerAPIKey, &mcpauth.RequireBearerTokenOptions{
		AllowMissingExpiration: true,
	})(streamable)
}

func newServer() *mcpsdk.Server {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    "macpanel",
		Version: global.CONF.Base.Version,
	}, nil)
	RegisterTools(server)
	return server
}

// RegisterRoutes mounts the MCP streamable HTTP endpoint on the Gin engine.
func RegisterRoutes(router *gin.Engine) {
	router.Any("/mcp", gin.WrapH(handler))
}
