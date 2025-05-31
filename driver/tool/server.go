package tool

import (
	"context"
	"errors"
)

var (
	ErrHandlerAlreadyExists = errors.New("handler already exists for this method")
	ErrMethodNotFound       = errors.New("method not found")
)

type Handler func(ctx context.Context, data []byte) (result []byte, err error)

type Server interface {
	AddHandler(method string, handler Handler) error
	Listen(ctx context.Context)
}
