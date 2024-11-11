package main

import (
	"backend/config"
	"backend/urls"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(config.SetupCors())

	config.ConnectDatabase()

	urls.InitializeRoutes(router)

	router.Run(":8080")
}
