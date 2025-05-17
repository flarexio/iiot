package http

import (
	"github.com/gin-gonic/gin"

	"github.com/flarexio/iiot"
)

func AddRouters(r *gin.Engine, endpoints iiot.EndpointSet) {
	// POST /iiot/check_connection
	r.POST("/iiot/check_connection", CheckConnectionHandler(endpoints.CheckConnection))
}
