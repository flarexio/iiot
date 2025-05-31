package http

import (
	"github.com/gin-gonic/gin"

	"github.com/flarexio/iiot"
)

func AddRouters(r *gin.Engine, endpoints iiot.EndpointSet) {
	r.POST("/iiot/check_connection", CheckConnectionHandler(endpoints.CheckConnection))
	r.GET("/iiot/drivers", ListDriversHandler(endpoints.ListDrivers))
	r.GET("/iiot/drivers/:driver/schema", SchemaHandler(endpoints.Schema))
	r.POST("/iiot/drivers/:driver/read_points", ReadPointsHandler(endpoints.ReadPoints))
}
