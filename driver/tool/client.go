package tool

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
	//   - schema: The schema in JSON format.
	//   - err: nil if the operation is successful, otherwise an error.
	Schema(ctx context.Context, driver string) (schema json.RawMessage, err error)

	// ReadPoints reads points from the given driver using the provided request.
	//
	// Args:
	//   - driver: The driver to use for reading points.
	//   - raw: The request in JSON format to be sent to the driver.
	// Returns:
	//   - results: A slice of results read from the driver, where each result can be of any type.
	//   - err: nil if the operation is successful, otherwise an error.
	ReadPoints(ctx context.Context, driver string, raw json.RawMessage) (results []any, err error)
}
