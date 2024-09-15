package main

import (
	"event-booking-api/app/router"
	"event-booking-api/config"
	_ "event-booking-api/docs"
	"log"
	"os"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	config.InitLog()
}

//	@title			event-booking-api swagger doc
//	@version		1.0
//	@description	event-booking-api swagger doc

//	@contact.name	devonchou
//	@contact.url	https://github.com/devonchou

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	port := os.Getenv("PORT")

	init := config.Init()
	app := router.Init(init)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := app.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
