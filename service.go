package iiot

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/flarexio/iiot/driver/tool"
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

	// ListDrivers retrieves a list of available drivers.
	//
	// Returns:
	//   - drivers: A slice of driver names.
	//   - error: nil if the operation is successful, otherwise an error.
	ListDrivers(ctx context.Context) (drivers []string, err error)

	tool.Client
}

type ServiceMiddleware func(Service) Service

func NewService(path string, tool tool.Client) Service {
	return &service{path, tool}
}

type service struct {
	path string
	tool tool.Client
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

func (svc *service) ListDrivers(ctx context.Context) ([]string, error) {
	driverPath := filepath.Join(svc.path, "drivers")

	entries, err := os.ReadDir(driverPath)
	if err != nil {
		return nil, err
	}

	drivers := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		driver, ok := strings.CutSuffix(filename, "_tool")
		if !ok {
			continue
		}

		drivers = append(drivers, driver)
	}

	if len(drivers) == 0 {
		return nil, errors.New("no drivers available")
	}

	return drivers, nil
}

func (svc *service) Schema(ctx context.Context, driver string) (json.RawMessage, error) {
	if driver == "" {
		return nil, errors.New("driver parameter is required")
	}

	return svc.tool.Schema(ctx, driver)
}

func (svc *service) Instruction(ctx context.Context, driver string) (string, error) {
	if driver == "" {
		return "", errors.New("driver parameter is required")
	}

	instruction, err := svc.tool.Instruction(ctx, driver)
	if err != nil {
		return "", err
	}

	if instruction == "" {
		return "", errors.New("no instruction available for the specified driver")
	}

	return instruction, nil
}

func (svc *service) ReadPoints(ctx context.Context, driver string, raw json.RawMessage) ([]any, error) {
	if driver == "" {
		return nil, errors.New("driver parameter is required")
	}

	return svc.tool.ReadPoints(ctx, driver, raw)
}
