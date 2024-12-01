package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	if os.Getenv("MODE") != "LIVE" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: No .env file found, using environment variables directly,\n")
		}
	}

	mode := os.Getenv("MODE")
	if mode == "" {
		log.Fatalf("Mode is not set in .env file")
	}

	dbNameKey := fmt.Sprintf("%s_POSTGRES_NAME", mode)
	dbUserKey := fmt.Sprintf("%s_POSTGRES_USER", mode)
	dbPasswordKey := fmt.Sprintf("%s_POSTGRES_PASSWORD", mode)
	dbHostKey := fmt.Sprintf("%s_POSTGRES_HOST", mode)
	dbPortKey := fmt.Sprintf("%s_POSTGRES_PORT", mode)

	dbName := os.Getenv(dbNameKey)
	dbUser := os.Getenv(dbUserKey)
	dbPassword := os.Getenv(dbPasswordKey)
	dbHost := os.Getenv(dbHostKey)
	dbPort := os.Getenv(dbPortKey)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}

	fmt.Println("Connected to the database")
	DB = pool
}

func SetupCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173", "https://portfolio-frontend-jesse-metzger.fly.dev", "https://jesse-metzger.com/",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// func SetupCors() gin.HandlerFunc {
// 	return cors.New(cors.Config{
// 		AllowOrigins: []string{
// 			"*",
// 		},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	})
// }
