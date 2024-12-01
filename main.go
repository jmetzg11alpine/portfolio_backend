package main

import (
	"backend/config"
	"backend/services/alpaca_script"
	"backend/urls"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
	router.Use(config.SetupCors())

	config.ConnectDatabase()

	urls.InitializeRoutes(router)

	c := cron.New()
	_, err = c.AddFunc("0 11 * * 1-5", func() {
		log.Println("Running scheduled alpaca script...")
		err := alpaca_script.Run()
		if err != nil {
			log.Printf("Failed to run alpaca script: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to schedule cron job: %v", err)
	}
	c.Start()

	router.Run(":8080")
}
