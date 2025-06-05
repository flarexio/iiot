package main

import (
	"context"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	"github.com/urfave/cli/v3"

	"github.com/flarexio/iiot"
	"github.com/flarexio/iiot/transport/mcp"
	"github.com/flarexio/iiot/transport/pubsub"
)

const (
	Version = "1.0.0"
)

func main() {
	cmd := &cli.Command{
		Name:  "iiot_mcp",
		Usage: "IIoT MCP server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "nats",
				Usage:   "NATS server URL",
				Value:   "wss://nats.flarex.io",
				Sources: cli.EnvVars("NATS_URL"),
			},
			&cli.StringFlag{
				Name:    "creds",
				Usage:   "NATS user credentials file",
				Sources: cli.EnvVars("NATS_CREDS"),
			},
		},
		Action: run,
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	natsURL := cmd.String("nats")
	natsCreds := cmd.String("creds")

	nc, err := nats.Connect(natsURL,
		nats.Name("IIoT MCP Server"),
		nats.UserCredentials(natsCreds),
	)

	if err != nil {
		return err
	}
	defer nc.Drain()

	topic := "edges.:edge_id.iiot"
	endpoints := pubsub.MakeEndpoints(nc, topic)

	// Create a new IIoT service
	var svc iiot.Service
	svc = iiot.ProxyMiddleware(endpoints)(svc)

	s := server.NewMCPServer(
		"IIoT Service",
		Version,
		server.WithToolHandlerMiddleware(mcp.InjectContextMiddleware()),
	)

	// Add CheckConnection tool
	{
		endpoint := iiot.CheckConnectionEndpoint(svc)
		handler := mcp.CheckConnectionHandler(endpoint)
		tool := mcp.CheckConnectionTool()
		s.AddTool(tool, handler)
	}

	// Add ListDrivers tool
	{
		endpoint := iiot.ListDriversEndpoint(svc)
		handler := mcp.ListDriversHandler(endpoint)
		tool := mcp.ListDriversTool()
		s.AddTool(tool, handler)
	}

	// Add Schema tool
	{
		endpoint := iiot.SchemaEndpoint(svc)
		handler := mcp.SchemaHandler(endpoint)
		tool := mcp.SchemaTool()
		s.AddTool(tool, handler)
	}

	// Add Instruction tool
	{
		endpoint := iiot.InstructionEndpoint(svc)
		handler := mcp.InstructionHandler(endpoint)
		tool := mcp.InstructionTool()
		s.AddTool(tool, handler)
	}

	// Add ReadPoints tool
	{
		endpoint := iiot.ReadPointsEndpoint(svc)
		handler := mcp.ReadPointsHandler(endpoint)
		tool := mcp.ReadPointsTool()
		s.AddTool(tool, handler)
	}

	return server.ServeStdio(s)
}
