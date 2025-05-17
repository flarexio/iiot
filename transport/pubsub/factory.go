package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/nats-io/nats.go"

	"github.com/flarexio/core/model"
	"github.com/flarexio/iiot"
)

func MakeEndpoints(nc *nats.Conn, prefix string) *iiot.EndpointSet {
	return &iiot.EndpointSet{
		CheckConnection: CheckConnectionEndpoint(nc, prefix+".check_connection"),
	}
}

func CheckConnectionEndpoint(nc *nats.Conn, topic string) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		if strings.Contains(topic, ":edge_id") {
			edgeID, ok := ctx.Value(model.EdgeID).(string)
			if !ok {
				return nil, errors.New("invalid edge id")
			}

			topic = strings.Replace(topic, ":edge_id", edgeID, 1)
		}

		req, ok := request.(iiot.CheckConnectionRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		data, err := json.Marshal(&req)
		if err != nil {
			return nil, err
		}

		msg, err := nc.Request(topic, data, 10*time.Second)
		if err != nil {
			return nil, err
		}

		return string(msg.Data), nil
	}
}
