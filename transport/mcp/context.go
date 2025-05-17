package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/flarexio/core/model"
)

type Property struct {
	name   string
	schema map[string]any
}

func NewProperty(name string, dataType string, opts ...mcp.PropertyOption) *Property {
	prop := &Property{
		name: name,
		schema: map[string]any{
			"type": dataType,
		},
	}

	for _, opt := range opts {
		opt(prop.schema)
	}
	delete(prop.schema, "required")

	return prop
}

func WithContext(name string, desc string, props ...*Property) mcp.ToolOption {
	properties := make(map[string]any)
	for _, prop := range props {
		properties[prop.name] = prop.schema
	}

	return mcp.WithObject(name,
		mcp.Description(desc),
		mcp.Properties(properties),
	)
}

func InjectContextMiddleware() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			c, ok := request.Params.Arguments["ctx"].(map[string]any)
			if !ok {
				return next(ctx, request)
			}

			edgeID, ok := c["edge_id"].(string)
			if ok {
				ctx = context.WithValue(ctx, model.EdgeID, edgeID)
			}

			return next(ctx, request)
		}
	}
}
