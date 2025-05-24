package pubsub

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	"github.com/nats-io/nats.go/micro"

	"github.com/flarexio/iiot"
)

func SchemaHandler(endpoint endpoint.Endpoint) micro.HandlerFunc {
	return func(r micro.Request) {
		driver := string(r.Data())
		if driver == "" {
			r.Error("400", "driver parameter is required", nil)
			return
		}

		ctx := context.Background()
		schema, err := endpoint(ctx, driver)
		if err != nil {
			r.Error("417", err.Error(), nil)
			return
		}

		r.RespondJSON(&schema)
	}
}

func ReadPointsHandler(endpoint endpoint.Endpoint) micro.HandlerFunc {
	return func(r micro.Request) {
		var req iiot.ReadPointsRequest
		if err := json.Unmarshal(r.Data(), &req); err != nil {
			r.Error("400", err.Error(), nil)
			return
		}

		ctx := context.Background()
		points, err := endpoint(ctx, req)
		if err != nil {
			r.Error("417", err.Error(), nil)
			return
		}

		r.RespondJSON(&points)
	}
}

func CheckConnectionHandler(endpoint endpoint.Endpoint) micro.HandlerFunc {
	return func(r micro.Request) {
		var req iiot.CheckConnectionRequest
		if err := json.Unmarshal(r.Data(), &req); err != nil {
			r.Error("400", err.Error(), nil)
			return
		}

		ctx := context.Background()
		_, err := endpoint(ctx, req)
		if err != nil {
			r.Error("417", err.Error(), nil)
			return
		}

		r.Respond([]byte("Connection successful"))
	}
}
