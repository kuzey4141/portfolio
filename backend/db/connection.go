package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool" // Use connection pool
)

var Pool *pgxpool.Pool // Use pool instead of single connection

// ConnectDB function establishes the connection with pool
func ConnectDB() {
	connStr := "postgresql://postgres:6303523aA@portfolyo.cdcik488ur35.eu-central-1.rds.amazonaws.com:5432/portfolio"

	log.Println("Attempting to connect to database...")
	log.Printf("Connection string: %s", connStr)

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
