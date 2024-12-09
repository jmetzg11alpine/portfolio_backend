package main

import (
	"backend/config"
	"backend/urls"
	"log"

	"github.com/gin-gonic/gin"
)

// func main() {
// 	gin.SetMode(gin.ReleaseMode)

// 	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
// 	if err != nil {
// 		log.Fatalf("Failed to open log file: %v", err)
// 	}

// 	log.SetOutput(logFile)

// 	log.Println("Go app started")

// 	router := gin.Default()
// 	if err := router.SetTrustedProxies(nil); err != nil {
// 		log.Fatalf("failed to set trusted proxies: %v", err)
// 	}
// 	router.Use(config.SetupCors())

// 	config.ConnectDatabase()

// 	urls.InitializeRoutes(router)

// 	log.Println("Adding cron job to run at 11am ET, Monday to Friday")
// 	c := cron.New()
// 	// second, minute, hour, day, month, day of week
// 	// sceduled to run 11AM ET, Monday - Friday
// 	err = c.AddFunc("0 0 11 * * 1-5", func() {
// 		log.Println("Running scheduled alpaca script...")
// 		err := alpaca_script.Run()
// 		if err != nil {
// 			log.Printf("Failed to run alpaca script: %v", err)
// 		}
// 	})
// 	if err != nil {
// 		log.Printf("Failed to schedule cron job: %v", err)
// 	}
// 	c.Start()

// 	router.Run(":8080")
// }

func main() {

	log.Println("Go app started")

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
	router.Use(config.SetupCors())

	config.ConnectDatabase()

	urls.InitializeRoutes(router)

	router.Run(":8080")
}
