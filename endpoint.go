package iiot

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type EndpointSet struct {
	Schema          endpoint.Endpoint
	ReadPoints      endpoint.Endpoint
	CheckConnection endpoint.Endpoint
}

func SchemaEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		driver, ok := request.(string)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.Schema(ctx, driver)
	}
}

type ReadPointsRequest struct {
	Driver string          `json:"driver"`
	Raw    json.RawMessage `json:"raw"`
}

func ReadPointsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(ReadPointsRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.ReadPoints(ctx, req.Driver, req.Raw)
	}
}

type CheckConnectionRequest struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

func CheckConnectionEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(CheckConnectionRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		err := svc.CheckConnection(ctx, req.Network, req.Address)
		return nil, err
	}
}
