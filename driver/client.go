package driver

import (
	"context"
	"encoding/json"
)

type Client interface {
	// Schema retrieves the schema for the specified driver.
	//
	// Args:
	//   - driver: The driver for which to retrieve the schema.
	// Returns:
	//   - json.RawMessage: The schema in JSON format.
	Schema(ctx context.Context, driver string) (json.RawMessage, error)

	// ReadPoints reads points from the given driver using the provided request.
	//
	// Args:
	//   - driver: The driver to use for reading points.
	//   - raw: The request in JSON format to be sent to the driver.
	// Returns:
	//   - []any: A slice of points read from the driver.
	//   - error: nil if the operation is successful, otherwise an error.
	ReadPoints(ctx context.Context, driver string, raw json.RawMessage) ([]any, error)
}
