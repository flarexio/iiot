package iiot

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type EndpointSet struct {
	CheckConnection endpoint.Endpoint
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
