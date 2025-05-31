package iiot

import (
	"context"
	"encoding/json"
	"errors"
)

func ProxyMiddleware(endpoints *EndpointSet) ServiceMiddleware {
	return func(next Service) Service {
		return &proxyMiddleware{
			endpoints: endpoints,
		}
	}
}

type proxyMiddleware struct {
	endpoints *EndpointSet
}

func (mw *proxyMiddleware) CheckConnection(ctx context.Context, network string, address string) error {
	req := CheckConnectionRequest{
		Network: network,
		Address: address,
	}

	_, err := mw.endpoints.CheckConnection(ctx, req)
	return err
}

func (mw *proxyMiddleware) ListDrivers(ctx context.Context) ([]string, error) {
	resp, err := mw.endpoints.ListDrivers(ctx, nil)
	if err != nil {
		return nil, err
	}

	drivers, ok := resp.([]string)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return drivers, nil
}

func (mw *proxyMiddleware) Schema(ctx context.Context, driver string) (json.RawMessage, error) {
	resp, err := mw.endpoints.Schema(ctx, driver)
	if err != nil {
		return nil, err
	}

	schema, ok := resp.(json.RawMessage)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return schema, nil
}

func (mw *proxyMiddleware) ReadPoints(ctx context.Context, driver string, raw json.RawMessage) ([]any, error) {
	req := ReadPointsRequest{
		Driver: driver,
		Raw:    raw,
	}

	resp, err := mw.endpoints.ReadPoints(ctx, req)
	if err != nil {
		return nil, err
	}

	points, ok := resp.([]any)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return points, nil
}
