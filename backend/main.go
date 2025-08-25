package main

import (
	"log"
	"os"
	"portfolio/db"
	"portfolio/routes"
	"portfolio/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	server.StartFrontend()

	db.ConnectDB()
	defer db.Pool.Close()

	r := gin.Default()
	routes.SetupRoutes(r, db.Pool)

	server.SetupStaticFiles(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on port %s", port)
	if gin.Mode() == gin.DebugMode {
		log.Println("Backend: http://localhost:" + port)
		log.Println("Frontend: http://localhost:3000")
	} else {
		log.Println("Application: http://localhost:" + port)
	}

	r.Run(":" + port)
}
