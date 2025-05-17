package iiot

import (
	"context"
	"net"
)

type Service interface {
	// CheckConnection checks if the given network and address are reachable.
	//
	// Args:
	//   - network: The network type (e.g., "tcp", "udp").
	//   - address: The address to check (e.g., "localhost:8080").
	// Returns:
	//   - error: nil if the connection is successful, otherwise an error.
	CheckConnection(ctx context.Context, network string, address string) error
}

type ServiceMiddleware func(Service) Service

func NewService() Service {
	return &service{}
}

type service struct {
}

func (svc *service) CheckConnection(ctx context.Context, network string, address string) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, network, address)
	if err != nil {
		return err
	}

	defer conn.Close()

	return nil
}
