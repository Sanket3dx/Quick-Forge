package main

import (
	"log"
	demonRoutes "quick_forge/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	demonRoutes.InitDemonRouter()
}
