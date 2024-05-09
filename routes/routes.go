package routes

import (
	service "goMqtt/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/api/grafana/meterData", func(c *gin.Context) {
		c.JSON(http.StatusOK, service.OneMinDataMap)
	})

	return router
}
