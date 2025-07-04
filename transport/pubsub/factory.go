package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/flarexio/core/model"
	"github.com/flarexio/iiot"
)

func MakeEndpoints(nc *nats.Conn, prefix string) *iiot.EndpointSet {
	return &iiot.EndpointSet{
		CheckConnection: CheckConnectionEndpoint(nc, prefix+".check_connection"),
		ListDrivers:     ListDriversEndpoint(nc, prefix+".drivers"),
		Schema:          SchemaEndpoint(nc, prefix+".schema"),
		Instruction:     InstructionEndpoint(nc, prefix+".instruction"),
		ReadPoints:      ReadPointsEndpoint(nc, prefix+".read_points"),
	}
}

func CheckConnectionEndpoint(nc *nats.Conn, topic string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		pubTopic := topic
		if strings.Contains(pubTopic, ":edge_id") {
			edgeID, ok := ctx.Value(model.EdgeID).(string)
			if !ok {
				return nil, errors.New("invalid edge id")
			}

			pubTopic = strings.Replace(pubTopic, ":edge_id", edgeID, 1)
		}

		req, ok := request.(iiot.CheckConnectionRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		data, err := json.Marshal(&req)
		if err != nil {
			return nil, err
		}

		msg, err := nc.Request(pubTopic, data, 10*time.Second)
		if err != nil {
			return nil, err
		}

		if err := Error(msg); err != nil {
			return nil, err
		}

		return string(msg.Data), nil
	}
}

func ListDriversEndpoint(nc *nats.Conn, topic string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		pubTopic := topic
		if strings.Contains(pubTopic, ":edge_id") {
			edgeID, ok := ctx.Value(model.EdgeID).(string)
			if !ok {
				return nil, errors.New("invalid edge id")
			}

			pubTopic = strings.Replace(pubTopic, ":edge_id", edgeID, 1)
		}

		msg, err := nc.Request(pubTopic, nil, nats.DefaultTimeout)
		if err != nil {
			return nil, err
		}

		if err := Error(msg); err != nil {
			return nil, err
		}

		var drivers []string
		if err := json.Unmarshal(msg.Data, &drivers); err != nil {
			return nil, err
		}

		return drivers, nil
	}
}

func SchemaEndpoint(nc *nats.Conn, topic string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		pubTopic := topic
		if strings.Contains(pubTopic, ":edge_id") {
			edgeID, ok := ctx.Value(model.EdgeID).(string)
			if !ok {
				return nil, errors.New("invalid edge id")
			}

			pubTopic = strings.Replace(pubTopic, ":edge_id", edgeID, 1)
		}

		driver, ok := request.(string)
		if !ok {
			return nil, errors.New("invalid request")
		}

		msg, err := nc.Request(pubTopic, []byte(driver), nats.DefaultTimeout)
		if err != nil {
			return nil, err
		}

		if err := Error(msg); err != nil {
			return nil, err
		}

		var schema json.RawMessage
		err = json.Unmarshal(msg.Data, &schema)
		if err != nil {
			return nil, err
		}

		return schema, nil
	}
}

func InstructionEndpoint(nc *nats.Conn, topic string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		pubTopic := topic
		if strings.Contains(pubTopic, ":edge_id") {
			edgeID, ok := ctx.Value(model.EdgeID).(string)
			if !ok {
				return nil, errors.New("invalid edge id")
			}

			pubTopic = strings.Replace(pubTopic, ":edge_id", edgeID, 1)
		}

		driver, ok := request.(string)
		if !ok {
			return nil, errors.New("invalid request")
		}

		msg, err := nc.Request(pubTopic, []byte(driver), nats.DefaultTimeout)
		if err != nil {
			return nil, err
		}

		if err := Error(msg); err != nil {
			return nil, err
		}

		return string(msg.Data), nil
	}
}

func ReadPointsEndpoint(nc *nats.Conn, topic string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		pubTopic := topic
		if strings.Contains(pubTopic, ":edge_id") {
			edgeID, ok := ctx.Value(model.EdgeID).(string)
			if !ok {
				return nil, errors.New("invalid edge id")
			}

			pubTopic = strings.Replace(pubTopic, ":edge_id", edgeID, 1)
		}

		req, ok := request.(iiot.ReadPointsRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		data, err := json.Marshal(&req)
		if err != nil {
			return nil, err
		}

		msg, err := nc.Request(pubTopic, data, nats.DefaultTimeout)
		if err != nil {
			return nil, err
		}

		if err := Error(msg); err != nil {
			return nil, err
		}

		var points []any
		err = json.Unmarshal(msg.Data, &points)
		if err != nil {
			return nil, err
		}

		return points, nil
	}
}

func Error(msg *nats.Msg) error {
	if msg == nil {
		return errors.New("nil message")
	}

	code := msg.Header.Get(micro.ErrorCodeHeader)
	if code == "" {
		return nil
	}

	description := msg.Header.Get(micro.ErrorHeader)
	if description == "" {
		description = "unknown error"
	}

	return errors.New(code + ": " + description)
}
