package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool" // Use connection pool
)

var Pool *pgxpool.Pool // Use pool instead of single connection

// ConnectDB function establishes the connection with pool
func ConnectDB() {
	connStr := "postgres://postgres:046804@localhost:5432/portfolio?sslmode=disable&pool_max_conns=10&pool_min_conns=2" // Pool settings added

	var err error
	Pool, err = pgxpool.New(context.Background(), connStr) // Use connection pool
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Test the connection
	if err := Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connected successfully with connection pool!")
}
