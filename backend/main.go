package main

import (
	"context"
	"portfolio/db"
	"portfolio/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()

	defer db.Conn.Close(context.Background())

	r := gin.Default()

	routes.SetupRoutes(r, db.Conn)

	r.Run(":8081") // Run on port 8081
}
