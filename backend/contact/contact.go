package contact

import (
	"context"        // For context in database operations
	"fmt"            // For printing and formatting to console
	"net/http"       // For HTTP status codes
	"portfolio/mail" // Mail package import
	"strconv"        // For string-integer conversion

	"github.com/gin-gonic/gin"        // Gin framework usage
	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL connection pool library
)

// Contact struct represents the contact table
type Contact struct {
	ID      int    `json:"id"`      // ID field, sent as "id" in JSON
	Name    string `json:"name"`    // Name field, sent as "name" in JSON
	Email   string `json:"email"`   // Email field, sent as "email" in JSON
	Phone   string `json:"phone"`   // Phone number, sent as "phone" in JSON
	Message string `json:"message"` // Message content, sent as "message" in JSON
}

var Pool *pgxpool.Pool // Global database pool is stored

func SetDB(pool *pgxpool.Pool) {
	Pool = pool // The incoming pool is assigned to the global Pool variable
}

func DeleteContact(c *gin.Context) {
	idStr := c.Param("id")         // Get ID from URL parameter (/api/contact/:id)
	id, err := strconv.Atoi(idStr) // Convert string ID to integer
	if err != nil {
		fmt.Println("ID conversion error:", err)                    // Print error to console
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"}) // Return 400 Bad Request
		return
	}

	_, err = Pool.Exec(context.Background(), "DELETE FROM contact WHERE id=$1", id) // Execute delete query
	if err != nil {
		fmt.Println("Delete error:", err)                                                 // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete operation failed"}) // Return 500 Internal Server Error
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Contact ID %d deleted successfully", id)}) // Return success message as JSON
}

func GetContacts(c *gin.Context) {
	rows, err := Pool.Query(context.Background(), "SELECT id, name, email, phone, message FROM contact") // Fetch all contact records
	if err != nil {
		fmt.Println("Data fetch error:", err)                                                 // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data could not be retrieved"}) // Return 500
		return
	}
	defer rows.Close() // Close resource at the end of the function

	var contacts []Contact // List of Contact structs
	for rows.Next() {      // Loop row by row
		var cct Contact
		if err := rows.Scan(&cct.ID, &cct.Name, &cct.Email, &cct.Phone, &cct.Message); err != nil {
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

	_, err := Pool.Exec(context.Background(),
		"UPDATE contact SET name=$1, email=$2, phone=$3, message=$4 WHERE id=$5",
		contact.Name, contact.Email, contact.Phone, contact.Message, contact.ID) // Execute update query
	if err != nil {
		fmt.Println("Database update error:", err)                              // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"}) // Return 500
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Contact ID %d updated successfully", contact.ID)}) // Return success message
}

// CreateContact - UPDATED with mail integration
func CreateContact(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil { // Bind JSON data to struct
		fmt.Println("JSON parsing error:", err)                       // Print error to console
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) // Return 400
		return
	}

	// Insert into database first
	_, err := Pool.Exec(context.Background(),
		"INSERT INTO contact (name, email, phone, message) VALUES ($1, $2, $3, $4)",
		contact.Name, contact.Email, contact.Phone, contact.Message) // Execute insert query
	if err != nil {
		fmt.Println("Database insert error:", err)                                          // Print error to console
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record could not be added"}) // Return 500
		return
	}

	// Send email notification
	mailData := mail.ContactMailData{
		Name:    contact.Name,
		Email:   contact.Email,
		Phone:   contact.Phone,
		Message: contact.Message,
	}

	// Send mail (error handling but don't block the response)
	if mailErr := mail.SendContactMail(mailData); mailErr != nil {
		fmt.Printf("Mail sending failed: %v\n", mailErr)
		// Even if mail fails, contact was saved, return success response
	} else {
		fmt.Println("Contact mail sent successfully!")
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Contact record added successfully and email sent"}) // Return success message with 201
}
