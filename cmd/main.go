package main

import (
	"courseworker/config"
	"courseworker/internal/handler"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbChan := make(chan *sql.DB, 1)
	rdcChan := make(chan *redis.Client, 1)
	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		db, err := config.SetupDB()
		if err != nil {
			errChan <- fmt.Errorf("database intialization error: %w", err)
			return
		}
		dbChan <- db
	}()

	go func() {
		defer wg.Done()
		redisClient, err := config.NewRedisClient()
		if err != nil {
			errChan <- fmt.Errorf("redis initialization error: %w", err)
			return
		}
		rdcChan <- redisClient
	}()

	wg.Wait()
	close(dbChan)
	close(rdcChan)
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			log.Printf("Error: %v", e)
		}
		log.Fatal("Initialization failed")
	}

	db := <-dbChan
	rdc := <-rdcChan
	defer db.Close()
	defer rdc.Close()

	r := gin.Default()
	handler.StartEngine(r, db, rdc)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
		log.Print("No APP_PORT found, using default: 8000")
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
