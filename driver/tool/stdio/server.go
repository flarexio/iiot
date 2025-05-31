package stdio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/flarexio/iiot/driver/tool"
)

var DefaultTimeout = 5000 * time.Millisecond

type StdioServer interface {
	tool.Server
	SetIO(in io.Reader, out io.Writer)
}

func NewStdioServer() StdioServer {
	return &stdioServer{
		in:       os.Stdin,
		out:      os.Stdout,
		handlers: make(map[string]tool.Handler),
	}
}

type stdioServer struct {
	in       io.Reader
	out      io.Writer
	handlers map[string]tool.Handler
	sync.RWMutex
}

func (s *stdioServer) AddHandler(method string, handler tool.Handler) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.handlers[method]; ok {
		return tool.ErrHandlerAlreadyExists
	}

	s.handlers[method] = handler
	return nil
}

func (s *stdioServer) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			err := s.handleRequest(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					select {
					case <-ctx.Done():
						return

					case <-time.After(100 * time.Millisecond):
						continue
					}
				}

				s.respond(nil, err)
			}
		}
	}
}

func (s *stdioServer) handleRequest(ctx context.Context) error {
	decoder := json.NewDecoder(s.in)

	var req *Request
	if err := decoder.Decode(&req); err != nil {
		return err
	}

	s.RLock()
	handler, ok := s.handlers[req.Method]
	s.RUnlock()

	if !ok {
		return tool.ErrMethodNotFound
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()

		result, err := handler(ctx, req.Data)
		s.respond(result, err)
	}()

	return nil
}

func (s *stdioServer) respond(result []byte, err error) error {
	resp := &Response{
		Result: result,
		Error:  err,
	}

	encoder := json.NewEncoder(s.out)
	return encoder.Encode(&resp)
}

func (s *stdioServer) SetIO(in io.Reader, out io.Writer) {
	s.in = in
	s.out = out
}
