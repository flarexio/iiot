package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/xeipuuv/gojsonschema"

	"github.com/flarexio/iiot/driver/tool"
	"github.com/flarexio/iiot/driver/tool/example"
	"github.com/flarexio/iiot/driver/tool/stdio"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tool := example.NewTool()

	server := stdio.NewStdioServer()
	server.AddHandler("driver.schema", SchemaHandler(tool))
	server.AddHandler("driver.instruction", InstructionHandler(tool))
	server.AddHandler("driver.readPoints", ReadPointsHandler(tool))

	go server.Listen(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}

func SchemaHandler(tool example.Tool) tool.Handler {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		return tool.Schema(ctx)
	}
}

func InstructionHandler(tool example.Tool) tool.Handler {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		instruction, err := tool.Instruction(ctx)
		if err != nil {
			return nil, err
		}

		return []byte(instruction), nil
	}
}

func ReadPointsHandler(tool example.Tool) tool.Handler {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		schema, err := tool.Schema(ctx)
		if err != nil {
			return nil, err
		}

		schemaLoader := gojsonschema.NewBytesLoader(schema)
		documentLoader := gojsonschema.NewBytesLoader(data)

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return nil, err
		}

		if !result.Valid() {
			var errs error
			for _, desc := range result.Errors() {
				errs = errors.Join(errs,
					fmt.Errorf("%s: %s", desc.Field(), desc.Description()))
			}

			return nil, errs
		}

		var req *example.ReadPointsRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}

		results, err := tool.ReadPoints(ctx, req)
		if err != nil {
			return nil, err
		}

		return json.Marshal(results)
	}
}
