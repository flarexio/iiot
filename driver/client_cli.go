package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
)

func NewCliDriverClient() Client {
	return &cliDriverClient{}
}

type cliDriverClient struct {
}

type Request struct {
	Method string
	Params map[string]any
}

type Response struct {
	Result any
	Error  error
}

func (resp *Response) UnmarshalJSON(data []byte) error {
	var raw struct {
		Result json.RawMessage
		Error  string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	resp.Result = raw.Result

	resp.Error = nil
	if raw.Error != "" {
		resp.Error = errors.New(raw.Error)
	}

	return nil
}

func (c *cliDriverClient) do(ctx context.Context, program string, req *Request) (*Response, error) {
	bs, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, program)

	var out bytes.Buffer
	cmd.Stdin = bytes.NewBuffer(bs)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var resp *Response
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *cliDriverClient) Schema(ctx context.Context, driver string) (json.RawMessage, error) {
	program := driver + "_cli"

	req := &Request{
		Method: "driver.schema",
		Params: nil,
	}

	resp, err := c.do(ctx, program, req)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	result, ok := resp.Result.(json.RawMessage)
	if !ok {
		return nil, errors.New("invalid response format")
	}

	return result, nil
}

func (c *cliDriverClient) ReadPoints(ctx context.Context, driver string, raw json.RawMessage) ([]any, error) {
	program := driver + "_cli"

	req := &Request{
		Method: "driver.readPoints",
		Params: map[string]any{
			"raw": raw,
		},
	}

	resp, err := c.do(ctx, program, req)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	results, ok := resp.Result.([]any)
	if !ok {
		return nil, errors.New("invalid response format")
	}

	return results, nil
}
