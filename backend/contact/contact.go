package contact // Defined as contact package

import (
	"context"  // For context in database operations
	"fmt"      // For printing and formatting to console
	"net/http" // For HTTP status codes
	"strconv"  // For string-integer conversion

	"github.com/gin-gonic/gin" // Gin framework usage
	"github.com/jackc/pgx/v5"  // PostgreSQL connection library
)

// Contact struct represents the contact table
type Contact struct {
	ID      int    `json:"id"`      // ID field, sent as "id" in JSON
	Email   string `json:"email"`   // Email field, sent as "email" in JSON
	Phone   string `json:"phone"`   // Phone number, sent as "phone" in JSON
	Message string `json:"message"` // Message content, sent as "message" in JSON
}

var Conn *pgx.Conn // Global database connection is stored

func SetDB(conn *pgx.Conn) {
	Conn = conn // The incoming connection is assigned to the global Conn variable
}

func DeleteContact(c *gin.Context) {
	idStr := c.Param("id")         // Get ID from URL parameter (/api/contact/:id)
	id, err := strconv.Atoi(idStr) // Convert string ID to integer
	if err != nil {
		fmt.Println("ID conversion error:", err)                    // Print error to console
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"}) // Return 400 Bad Request
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM contact WHERE id=$1", id) // Execute delete query
	if err != nil {
		fmt.Println("Delete error:", err)                                                 // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete operation failed"}) // Return 500 Internal Server Error
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Contact ID %d deleted successfully", id)}) // Return success message as JSON
}

func GetContacts(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, email, phone, message FROM contact") // Fetch all contact records
	if err != nil {
		fmt.Println("Data fetch error:", err)                                                 // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data could not be retrieved"}) // Return 500
		return
	}
	defer rows.Close() // Close resource at the end of the function

	var contacts []Contact // List of Contact structs
	for rows.Next() {      // Loop row by row
		var cct Contact
		if err := rows.Scan(&cct.ID, &cct.Email, &cct.Phone, &cct.Message); err != nil {
			fmt.Println("Error reading row:", err)                                          // Print error to console
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row could not be read"}) // Return 500
			return
		}
		contacts = append(contacts, cct) // Add to list
	}

	c.JSON(http.StatusOK, contacts) // Return all records as JSON
}

func UpdateContact(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil { // Bind JSON data to struct
		fmt.Println("JSON parsing error:", err)                       // Print error to console
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) // Return 400
		return
	}

	_, err := Conn.Exec(context.Background(),
		"UPDATE contact SET email=$1, phone=$2, message=$3 WHERE id=$4",
		contact.Email, contact.Phone, contact.Message, contact.ID) // Execute update query
	if err != nil {
		fmt.Println("Database update error:", err)                              // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"}) // Return 500
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Contact ID %d updated successfully", contact.ID)}) // Return success message
}

func CreateContact(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil { // Bind JSON data to struct
		fmt.Println("JSON parsing error:", err)                       // Print error to console
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) // Return 400
		return
	}

	_, err := Conn.Exec(context.Background(),
		"INSERT INTO contact (email, phone, message) VALUES ($1, $2, $3)",
		contact.Email, contact.Phone, contact.Message) // Execute insert query
	if err != nil {
		fmt.Println("Database insert error:", err)                                          // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record could not be added"}) // Return 500
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "New contact record added successfully"}) // Return success message with 201
}
