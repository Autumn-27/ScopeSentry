package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "scopesentry",
		Version: constants.Version,
	}, nil)

	registerTools(server)
	return server
}

func jsonToolResult(v any) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}

func errorResult(msg string, err error) (*mcp.CallToolResult, any, error) {
	text := msg
	if err != nil {
		text = fmt.Sprintf("%s: %v", msg, err)
	}
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, nil, nil
}
