package home

import (
	"context"  // Context for database operations
	"fmt"      // For formatted printing
	"net/http" // For HTTP status codes
	"strconv"  // For string to int conversion

	"github.com/gin-gonic/gin"        // Gin framework
	"github.com/jackc/pgx/v5/pgxpool" // pgxpool library for PostgreSQL connection pool
)

// Home struct represents the data in the home table
type Home struct {
	ID          int    `json:"id"`          // Appears as "id" in JSON output
	Title       string `json:"title"`       // Title field
	Description string `json:"description"` // Description field
}

// Pool is a global variable for the database connection pool
var Pool *pgxpool.Pool

// SetDB function receives the database pool from main.go
func SetDB(pool *pgxpool.Pool) {
	Pool = pool
}

// DeleteHome deletes a home record by a specific ID
func DeleteHome(c *gin.Context) {
	idStr := c.Param("id")         // Get id from URL parameter (/api/home/:id format)
	id, err := strconv.Atoi(idStr) // Convert string to integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"}) // If invalid ID, return 400
		return
	}

	_, err = Pool.Exec(context.Background(), "DELETE FROM home WHERE id=$1", id) // Delete query from database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete operation failed"}) // If delete fails, return 500
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Home ID %d deleted successfully", id)}) // Return success message
}

// GetHomes returns all records from the home table when an HTTP GET request is received
func GetHomes(c *gin.Context) {
	rows, err := Pool.Query(context.Background(), "SELECT id, title, description FROM home") // Select all home records
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data could not be retrieved"}) // If data cannot be fetched, return 500
		return
	}
	defer rows.Close() // Close the connection when done

	var homes []Home  // Create empty slice
	for rows.Next() { // Loop through rows
		var h Home
		if err := rows.Scan(&h.ID, &h.Title, &h.Description); err != nil { // Get data from row
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row could not be read"}) // If error, return 500
			return
		}
		homes = append(homes, h) // Add to slice
	}

	c.JSON(http.StatusOK, homes) // Return all records as JSON
}

// UpdateHome updates a home record
func UpdateHome(c *gin.Context) {
	var h Home
	if err := c.ShouldBindJSON(&h); err != nil { // Bind JSON data to struct (decode)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) // If JSON invalid, return 400
		return
	}

	_, err := Pool.Exec(context.Background(), "UPDATE home SET title=$1, description=$2 WHERE id=$3", h.Title, h.Description, h.ID) // Update query
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"}) // If fails, return 500
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Home ID %d updated successfully", h.ID)}) // Return success message
}

// CreateHome function adds a new home record
func CreateHome(c *gin.Context) {
	var h Home
	if err := c.ShouldBindJSON(&h); err != nil { // Bind JSON data to struct
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) // If invalid data, return 400
		return
	}

	_, err := Pool.Exec(context.Background(), "INSERT INTO home (title, description) VALUES ($1, $2)", h.Title, h.Description) // Insert new record
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record could not be added"}) // If cannot be added, return 500
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Home record added successfully"}) // Return success message (201)
}
