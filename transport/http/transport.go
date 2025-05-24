package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"

	"github.com/flarexio/iiot"
)

func SchemaHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		driver := c.Param("driver")
		if driver == "" {
			err := errors.New("driver parameter is required")
			c.String(http.StatusBadRequest, err.Error())
			c.Error(err)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		schema, err := endpoint(ctx, driver)
		if err != nil {
			c.String(http.StatusExpectationFailed, err.Error())
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, schema)
	}
}

func ReadPointsHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req iiot.ReadPointsRequest
		if err := c.ShouldBind(&req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Error(err)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		points, err := endpoint(ctx, req)
		if err != nil {
			c.String(http.StatusExpectationFailed, err.Error())
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, points)
	}
}

func CheckConnectionHandler(endpoint endpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req iiot.CheckConnectionRequest
		if err := c.ShouldBind(&req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Error(err)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		_, err := endpoint(ctx, req)
		if err != nil {
			c.String(http.StatusExpectationFailed, err.Error())
			c.Error(err)
			c.Abort()
			return
		}

		c.String(http.StatusOK, "Connection successful")
	}
}
