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

type EdgeContext struct {
	EdgeID string `json:"edge_id"`
}

func InjectContextMiddleware() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req struct {
				EdgeCtx EdgeContext `json:"ctx"`
			}

			if err := request.BindArguments(&req); err != nil {
				return next(ctx, request)
			}

			if req.EdgeCtx.EdgeID != "" {
				ctx = context.WithValue(ctx, model.EdgeID, req.EdgeCtx.EdgeID)
			}

			return next(ctx, request)
		}
	}
}
