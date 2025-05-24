package iiot

import (
	"context"
	"encoding/json"
	"net"

	"github.com/flarexio/iiot/driver"
)

type Service interface {
	driver.Client

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

func NewService(client driver.Client) Service {
	return &service{client}
}

type service struct {
	client driver.Client
}

func (svc *service) Schema(ctx context.Context, driver string) (json.RawMessage, error) {
	return svc.client.Schema(ctx, driver)
}

func (svc *service) ReadPoints(ctx context.Context, driver string, raw json.RawMessage) ([]any, error) {
	return svc.client.ReadPoints(ctx, driver, raw)
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
