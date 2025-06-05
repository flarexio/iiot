package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	"github.com/flarexio/iiot"
	"github.com/flarexio/iiot/transport/pubsub"

	mcptransport "github.com/flarexio/iiot/transport/mcp"
)

func TestMCPToolCall(t *testing.T) {
	assert := assert.New(t)

	natsURL := os.Getenv("NATS_URL")
	natsCreds := os.Getenv("NATS_CREDS")

	nc, err := nats.Connect(
		natsURL,
		nats.Name("IIoT MCP Server Test"),
		nats.UserCredentials(natsCreds),
	)

	if err != nil {
		assert.Fail(err.Error())
		return
	}
	defer nc.Drain()

	topic := "edges.:edge_id.iiot"
	endpoints := pubsub.MakeEndpoints(nc, topic)

	var svc iiot.Service
	svc = iiot.ProxyMiddleware(endpoints)(svc)

	s := server.NewMCPServer(
		"IIoT Service",
		"1.0.0",
		server.WithToolHandlerMiddleware(
			mcptransport.InjectContextMiddleware(),
		),
	)

	endpoint := iiot.ListDriversEndpoint(svc)
	handler := mcptransport.ListDriversHandler(endpoint)
	tool := mcptransport.ListDriversTool()
	s.AddTool(tool, handler)

	req := mcp.JSONRPCRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      mcp.NewRequestId(1),
		Params: mcp.Params{
			"name": "ListDrivers",
			"arguments": map[string]any{
				"ctx": map[string]any{
					"edge_id": "01J6TRZ0RWW334GPRMH5NSKJQA",
				},
			},
		},
		Request: mcp.Request{
			Method: "tools/call",
		},
	}

	bs, err := json.Marshal(&req)
	if err != nil {
		assert.Fail(err.Error())
		return
	}

	ctx := context.Background()
	respRpc := s.HandleMessage(ctx, bs)

	resp, ok := respRpc.(mcp.JSONRPCResponse)
	if !ok {
		assert.Fail("invalid response type")
		return
	}

	result, ok := resp.Result.(mcp.CallToolResult)
	if !ok {
		assert.Fail("invalid result type")
		return
	}

	assert.Equal(mcp.JSONRPC_VERSION, resp.JSONRPC)
	assert.Equal(mcp.NewRequestId(float64(1)), resp.ID)
	assert.False(result.IsError)
	assert.Len(result.Content, 1)

	content, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		assert.Fail("invalid content type")
		return
	}

	var drivers []string
	if err := json.Unmarshal([]byte(content.Text), &drivers); err != nil {
		assert.Fail(err.Error())
		return
	}

	assert.Contains(drivers, "example")
	assert.Contains(drivers, "modbus")
}
