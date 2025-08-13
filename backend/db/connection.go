// backend/db/connection.go
package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5" // PostgreSQL library
)

var Conn *pgx.Conn // Global database connection

// ConnectDB function establishes the connection
func ConnectDB() {
	connStr := "postgres://postgres:046804@localhost:5432/portfolio?sslmode=disable" // Connection string

	var err error
	Conn, err = pgx.Connect(context.Background(), connStr) // Connect to the database
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err) // Log error and exit if connection fails
	}
}
