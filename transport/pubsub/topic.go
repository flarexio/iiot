package pubsub

import (
	"github.com/nats-io/nats.go/micro"

	"github.com/flarexio/iiot"
)

func AddEndpoints(group micro.Group, endpoints iiot.EndpointSet) {
	group.AddEndpoint("check_connection", CheckConnectionHandler(endpoints.CheckConnection))
}
