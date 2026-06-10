package mcp

import (
	"net/http"

	"github.com/Autumn-27/ScopeSentry/internal/api/middleware"
	"github.com/gin-gonic/gin"
	mcpSDK "github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterRoutes 注册 MCP Streamable HTTP 路由
func RegisterRoutes(router *gin.Engine) {
	handler := mcpSDK.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcpSDK.Server { return newServer() },
		&mcpSDK.StreamableHTTPOptions{Stateless: true},
	)

	mcpGroup := router.Group("/mcp")
	mcpGroup.Use(middleware.MCPAPIKeyMiddleware())
	mcpGroup.Any("", gin.WrapH(handler))
	mcpGroup.Any("/*filepath", gin.WrapH(handler))
}
