package mcp

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/flarexio/iiot"
)

func CheckConnectionTool(name ...string) mcp.Tool {
	toolName := "check_connection"
	if len(name) > 0 {
		toolName = name[0]
	}

	return mcp.NewTool(toolName,
		mcp.WithDescription("Check the connection to a given network and address"),
		WithContext("ctx", "Context for the request",
			NewProperty("edge_id", "string", mcp.Description("The edge ID for the request")),
		),
		mcp.WithString("network",
			mcp.Required(),
			mcp.Description("The network type (e.g., tcp, udp)"),
			mcp.DefaultString("tcp"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The address to check (e.g., 192.168.1.100:502)"),
		),
	)
}

func CheckConnectionHandler(endpoint endpoint.Endpoint) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments

		network, ok := args["network"].(string)
		if !ok {
			return nil, errors.New("invalid network type")
		}

		address, ok := args["address"].(string)
		if !ok {
			return nil, errors.New("invalid address")
		}

		req := iiot.CheckConnectionRequest{
			Network: network,
			Address: address,
		}

		_, err := endpoint(ctx, req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText("Connection successful"), nil
	}
}
