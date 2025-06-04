package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/flarexio/iiot"
	"github.com/flarexio/iiot/transport/http"
	"github.com/flarexio/iiot/transport/pubsub"

	tool "github.com/flarexio/iiot/driver/tool/stdio"
)

const (
	Version = "1.0.0"
)

func main() {
	cmd := &cli.Command{
		Name:  "iiot",
		Usage: "IIoT Service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "Path to the IIoT service",
			},
			&cli.IntFlag{
				Name:  "port",
				Usage: "HTTP server port",
				Value: 8080,
			},
			&cli.StringFlag{
				Name:    "nats",
				Usage:   "NATS server URL",
				Value:   "wss://nats.flarex.io",
				Sources: cli.EnvVars("NATS_URL"),
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
	path := cmd.String("path")
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		path = homeDir + "/.flarex/iiot"
	}

	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer log.Sync()

	zap.ReplaceGlobals(log) // Replace the global logger

	// Initialize the tool client
	driverPath := filepath.Join(path, "drivers")
	executor := tool.NewCommandExecutor(driverPath)
	tool := tool.NewStdioClient(executor)

	// Create a new IIoT service
	svc := iiot.NewService(path, tool)
	svc = iiot.LoggingMiddleware(log)(svc)

	endpoints := iiot.EndpointSet{
		CheckConnection: iiot.CheckConnectionEndpoint(svc),
		ListDrivers:     iiot.ListDriversEndpoint(svc),
		Schema:          iiot.SchemaEndpoint(svc),
		Instruction:     iiot.InstructionEndpoint(svc),
		ReadPoints:      iiot.ReadPointsEndpoint(svc),
	}

	// Add HTTP Transport
	{
		port := cmd.Int("port")

		r := gin.Default()
		http.AddRouters(r, endpoints)

		go r.Run(":" + strconv.Itoa(port))
	}

	// Add PubSub Transport
	{
		natsURL := cmd.String("nats")
		natsCreds := filepath.Join(path, "user.creds")

		f, err := os.ReadFile(filepath.Join(path, "id"))
		if err != nil {
			return err
		}

		edgeID := strings.TrimSpace(string(f))

		nc, err := nats.Connect(natsURL,
			nats.Name("IIoT Service - "+edgeID),
			nats.UserCredentials(natsCreds),
		)

		if err != nil {
			return err
		}
		defer nc.Drain()

		srv, err := micro.AddService(nc, micro.Config{
			Name:    "iiot",
			Version: Version,
		})

		if err != nil {
			return err
		}
		defer srv.Stop()

		topic := "edges." + edgeID + ".iiot"

		root := srv.AddGroup(topic)
		pubsub.AddEndpoints(root, endpoints)
	}

	// Setup signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit // Wait for a termination signal

	log.Info("graceful shutdown", zap.String("signal", sign.String()))
	return nil
}
