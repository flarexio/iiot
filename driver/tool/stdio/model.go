package stdio

import (
	"encoding/json"
	"errors"
)

type Request struct {
	Method string
	Data   []byte
}

type Response struct {
	Result []byte
	Error  error
}

func (resp *Response) UnmarshalJSON(data []byte) error {
	var raw struct {
		Result []byte
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

func (resp *Response) MarshalJSON() ([]byte, error) {
	raw := struct {
		Result []byte `json:"result"`
		Error  string `json:"error,omitempty"`
	}{
		Result: resp.Result,
	}

	if resp.Error != nil {
		raw.Error = resp.Error.Error()
	}

	return json.Marshal(raw)
}
