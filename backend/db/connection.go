package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool" // Use connection pool
)

var Pool *pgxpool.Pool // Use pool instead of single connection

// ConnectDB function establishes the connection with pool
func ConnectDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	log.Println("Attempting to connect to database...")

	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Pool created, testing connection...")

	if err := Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connected successfully with connection pool!")
}
