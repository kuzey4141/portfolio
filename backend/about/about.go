package about // about package is defined, all code in this file belongs to this package

import (
	"context"  // context is used for database operations
	"fmt"      // For writing error or information messages to the terminal
	"net/http" // For HTTP status codes
	"strconv"  // To convert string expressions to integer (for example ID)

	"github.com/gin-gonic/gin" // To use the Gin framework
	"github.com/jackc/pgx/v5"  // To communicate with PostgreSQL database using pgx library
)

// About struct represents a row in the "about" table in the database
type About struct {
	ID      int    `json:"id"`      // Displayed as "id" field in JSON output
	Content string `json:"content"` // Displayed as "content" field in JSON output
}

var Conn *pgx.Conn // Variable to hold the database connection globally

// SetDB function sets the externally received database connection to be used in this package
func SetDB(conn *pgx.Conn) {
	Conn = conn // Incoming connection is assigned to the global variable
}

// GetAbouts function retrieves all "about" records from the database and returns them as JSON
func GetAbouts(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, content FROM about") // SQL query is executed
	if err != nil {
		fmt.Println(err)                                                                      // If there is an error, print to terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data could not be retrieved"}) // Return HTTP 500 error as JSON
		return
	}
	defer rows.Close() // Release resources after the query is completed

	var abouts []About // Empty slice to store incoming data
	for rows.Next() {  // Loop for each row
		var a About                                          // Temporary About object
		if err := rows.Scan(&a.ID, &a.Content); err != nil { // Row data is assigned to struct
			fmt.Println(err)                                                                // Print error if exists
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row could not be read"}) // Return error as JSON
			return
		}
		abouts = append(abouts, a) // Add struct to slice
	}

	c.JSON(http.StatusOK, abouts) // Return all records as JSON
}

// DeleteAbout function deletes the about record with the specified ID
func DeleteAbout(c *gin.Context) {
	idStr := c.Param("id")         // Get ID from URL parameter (/api/about/:id format)
	id, err := strconv.Atoi(idStr) // Convert string ID to integer
	if err != nil {
		fmt.Println(err)                                            // Print error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"}) // Return 400 if invalid ID
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM about WHERE id=$1", id) // Execute delete query
	if err != nil {
		fmt.Println(err)                                                                  // Print error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete operation failed"}) // Return 500 if failed
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "About record deleted."}) // Return success message as JSON
}

// UpdateAbout function updates the content of the about record with the specified ID
func UpdateAbout(c *gin.Context) {
	var a About
	if err := c.ShouldBindJSON(&a); err != nil { // Bind JSON data to struct
		fmt.Println("JSON parsing error:", err)                       // Print error if exists
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) // Return 400 if invalid data
		return
	}
	var ad string

	fmt.Println("What is your name?:")
	fmt.Scan(&ad)

	fmt.Println("Welcome,", ad)

	_, err := Conn.Exec(context.Background(), "UPDATE about SET content=$1 WHERE id=$2", a.Content, a.ID) // Execute update query
	if err != nil {
		fmt.Println("Database update error:", err)                              // Print error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"}) // Return 500 if failed
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("About ID %d updated successfully", a.ID)}) // Return success message as JSON
}

// CreateAbout function adds a new about record
func CreateAbout(c *gin.Context) {
	var a About
	if err := c.ShouldBindJSON(&a); err != nil { // Bind JSON data to struct
		fmt.Println("JSON parsing error:", err)                       // Print error if exists
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) // Return 400 if invalid data
		return
	}

	_, err := Conn.Exec(context.Background(), "INSERT INTO about (content) VALUES ($1)", a.Content) // Execute insert query
	if err != nil {
		fmt.Println("Database insert error:", err)                                          // Print error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record could not be added"}) // Return 500 if failed
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "New About record added successfully"}) // Return success message with 201
}
