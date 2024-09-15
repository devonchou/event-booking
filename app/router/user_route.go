package router

import (
	"event-booking-api/app/middleware"
	"event-booking-api/config"

	"github.com/gin-gonic/gin"
)

func addUserRoute(rg *gin.RouterGroup, init *config.Initialization) {
	user := rg.Group("/users")

	user.POST("", init.UserCtrl.AddUser)
	user.POST("/login", init.UserCtrl.LoginUser)

	protected := user.Group("")
	protected.Use(middleware.Auth)
	protected.GET("", init.UserCtrl.GetAllUser)
	protected.GET("/:userId", init.UserCtrl.GetUserById)
	protected.PUT("/:userId", init.UserCtrl.UpdateUserById)
	protected.DELETE("/:userId", init.UserCtrl.DeleteUserById)
}
