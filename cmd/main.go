package main

import (
	"courseworker/config"
	"courseworker/internal/handler"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.SetupDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	redisClient := config.NewRedisClient()
	defer redisClient.Close()

	r := gin.Default()
	handler.StartEngine(r, db, redisClient)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	r.Run(":" + port)
}
