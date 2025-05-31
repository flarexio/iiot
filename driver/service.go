package driver

import (
	"context"
	"errors"

	"github.com/flarexio/iiot/machine"
)

var (
	ErrControllerNotFound = errors.New("controller not found")
	ErrPointNotFound      = errors.New("point not found")
)

type Service interface {
	AddControllers(controllers ...*machine.Controller) error
	ReadPoints(ctx context.Context, id string, pointNames []string) (points []any, err error)
}
