package main

import (
	"log"
	"os"
	"portfolio/db"
	"portfolio/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println(".env file loaded successfully!")
	}

	log.Printf("RESEND_API_KEY length: %d", len(os.Getenv("RESEND_API_KEY")))
	log.Printf("TO_EMAIL: %s", os.Getenv("TO_EMAIL"))

	db.ConnectDB()
	defer db.Pool.Close() // Close pool when program exits

	r := gin.Default()
	routes.SetupRoutes(r, db.Pool)

	r.Run(":8081")
}
