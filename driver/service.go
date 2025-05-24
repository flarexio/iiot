package driver

import (
	"context"
	"encoding/json"
)

type Service interface {
	Schema(ctx context.Context) (json.RawMessage, error)
	ReadPoints(ctx context.Context, raw json.RawMessage) ([]any, error)
}
