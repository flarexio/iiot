package stdio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/flarexio/iiot/driver/tool"
)

type StdioClient interface {
	tool.Client
}

func NewStdioClient(executor Executor) StdioClient {
	return &stdioClient{executor}
}

type stdioClient struct {
	executor Executor
}

func (c *stdioClient) do(ctx context.Context, program string, req *Request) (*Response, error) {
	if c.executor == nil {
		return nil, errors.New("executor is not set")
	}

	in := new(bytes.Buffer)
	out := new(bytes.Buffer)

	encoder := json.NewEncoder(in)
	if err := encoder.Encode(&req); err != nil {
		return nil, err
	}

	if err := c.executor.Execute(ctx, program, in, out); err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(out)

	var resp *Response
	if err := decoder.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *stdioClient) Schema(ctx context.Context, driver string) (json.RawMessage, error) {
	program := driver + "_tool"

	req := &Request{
		Method: "driver.schema",
	}

	resp, err := c.do(ctx, program, req)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Result, nil
}

func (c *stdioClient) Instruction(ctx context.Context, driver string) (string, error) {
	program := driver + "_tool"

	req := &Request{
		Method: "driver.instruction",
	}

	resp, err := c.do(ctx, program, req)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}

	return string(resp.Result), nil
}

func (c *stdioClient) ReadPoints(ctx context.Context, driver string, raw json.RawMessage) ([]any, error) {
	program := driver + "_tool"

	req := &Request{
		Method: "driver.readPoints",
		Data:   raw,
	}

	resp, err := c.do(ctx, program, req)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	var results []any
	if err := json.Unmarshal(resp.Result, &results); err != nil {
		return nil, err
	}

	return results, nil
}
