package pubsub

import (
	"github.com/nats-io/nats.go/micro"

	"github.com/flarexio/iiot"
)

func AddEndpoints(group micro.Group, endpoints iiot.EndpointSet) {
	group.AddEndpoint("check_connection", CheckConnectionHandler(endpoints.CheckConnection))
	group.AddEndpoint("drivers", ListDriversHandler(endpoints.ListDrivers))
	group.AddEndpoint("schema", SchemaHandler(endpoints.Schema))
	group.AddEndpoint("read_points", ReadPointsHandler(endpoints.ReadPoints))
}
