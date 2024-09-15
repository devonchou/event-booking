package router

import (
	"event-booking-api/app/middleware"
	"event-booking-api/config"

	"github.com/gin-gonic/gin"
)

func addEventRoute(rg *gin.RouterGroup, init *config.Initialization) {
	event := rg.Group("/events")

	event.GET("", init.EventCtrl.GetAllEvent)
	event.GET("/:eventId", init.EventCtrl.GetEventById)

	protected := event.Group("")
	protected.Use(middleware.Auth)
	protected.POST("", init.EventCtrl.AddEvent)
	protected.PUT("/:eventId", init.EventCtrl.UpdateEventById)
	protected.DELETE("/:eventId", init.EventCtrl.DeleteEventById)
	protected.POST("/:eventId/register", init.EventCtrl.RegisterUserForEvent)
	protected.DELETE("/:eventId/register", init.EventCtrl.UnregisterUserForEvent)
	protected.GET("/:eventId/attendees", init.EventCtrl.GetAttendeesEmailById)
}
