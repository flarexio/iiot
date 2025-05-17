package iiot

import "context"

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
