package main

import (
	"courseworker/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.SetupDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
