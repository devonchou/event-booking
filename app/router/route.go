package router

import (
	"event-booking-api/config"

	"github.com/gin-gonic/gin"
)

func Init(init *config.Initialization) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	addUserRoute(api, init)
	addEventRoute(api, init)

	return router
}
