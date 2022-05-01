package main

import (
	"github.com/decadev/squad10/healthplus/db"
	"github.com/decadev/squad10/healthplus/router"
	"github.com/joho/godotenv"
	"log"
)

// start
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.SetupDB()

	router.SetupRouter()
}
