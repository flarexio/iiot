package mcp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/flarexio/iiot"
)

func CheckConnectionTool(name ...string) mcp.Tool {
	toolName := "CheckConnection"
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

func ListDriversTool(name ...string) mcp.Tool {
	toolName := "ListDrivers"
	if len(name) > 0 {
		toolName = name[0]
	}

	return mcp.NewTool(toolName,
		mcp.WithDescription("List all available drivers for IIoT protocols."),
		WithContext("ctx", "Context for the request",
			NewProperty("edge_id", "string", mcp.Description("The edge ID for the request")),
		),
	)
}

func ListDriversHandler(endpoint endpoint.Endpoint) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := endpoint(ctx, nil)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		drivers, ok := resp.([]string)
		if !ok {
			err = errors.New("invalid response type")
			return mcp.NewToolResultError(err.Error()), nil
		}

		bs, err := json.Marshal(&drivers)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(bs)), nil
	}
}

func SchemaTool(name ...string) mcp.Tool {
	toolName := "Schema"
	if len(name) > 0 {
		toolName = name[0]
	}

	return mcp.NewTool(toolName,
		mcp.WithDescription("Get the schema for a specific driver configuration."),
		WithContext("ctx", "Context for the request",
			NewProperty("edge_id", "string", mcp.Description("The edge ID for the request")),
		),
		mcp.WithString("driver",
			mcp.Required(),
			mcp.Description("The name of the driver, such as 'modbus', 'opcua', etc."),
		),
	)
}

func SchemaHandler(endpoint endpoint.Endpoint) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments

		driver, ok := args["driver"].(string)
		if !ok {
			return nil, errors.New("invalid driver type")
		}

		resp, err := endpoint(ctx, driver)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		schema, ok := resp.(json.RawMessage)
		if !ok {
			err := errors.New("invalid response type")
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(schema)), nil
	}
}

func ReadPointsTool(name ...string) mcp.Tool {
	toolName := "ReadPoints"
	if len(name) > 0 {
		toolName = name[0]
	}

	return mcp.NewTool(toolName,
		mcp.WithDescription("Read data points using a specific protocol driver and configuration. Use Schema tool first to get the required configuration format."),
		WithContext("ctx", "Context for the request",
			NewProperty("edge_id", "string", mcp.Description("The edge ID for the request")),
		),
		mcp.WithString("driver",
			mcp.Required(),
			mcp.Description("The name of the driver, such as 'modbus', 'opcua', etc."),
		),
		mcp.WithObject("raw",
			mcp.Required(),
			mcp.Description("Driver-specific configuration object. Use Schema to get the required format and structure."),
		),
	)
}

func ReadPointsHandler(endpoint endpoint.Endpoint) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments

		driver, ok := args["driver"].(string)
		if !ok {
			return nil, errors.New("invalid driver type")
		}

		raw, ok := args["raw"].(json.RawMessage)
		if !ok {
			return nil, errors.New("invalid raw data")
		}

		req := iiot.ReadPointsRequest{
			Driver: driver,
			Raw:    raw,
		}

		resp, err := endpoint(ctx, req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		points, ok := resp.([]any)
		if !ok {
			err = errors.New("invalid response type")
			return mcp.NewToolResultError(err.Error()), nil
		}

		bs, err := json.Marshal(&points)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(bs)), nil
	}
}
