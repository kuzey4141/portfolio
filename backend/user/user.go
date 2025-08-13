package user // user package

import (
	"context" // context for DB operations
	"fmt"     // for printing to console
	"strconv" // for string - int conversion

	"github.com/gin-gonic/gin" // Gin framework
	"github.com/jackc/pgx/v5"  // PostgreSQL library
)

// User struct represents data in the users table
type User struct {
	ID       int    `json:"id"`                 // id field in JSON
	Username string `json:"username"`           // username field in JSON
	Password string `json:"password,omitempty"` // password in JSON (shown empty if not provided)
	Email    string `json:"email"`              // email in JSON
}

var Conn *pgx.Conn // Global database connection

func SetDB(conn *pgx.Conn) {
	Conn = conn // Set the connection
}

// DeleteUser delete operation
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")         // Get id parameter from URL
	id, err := strconv.Atoi(idStr) // Convert string to int
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"}) // If invalid ID, return JSON error
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id) // Delete from DB
	if err != nil {
		c.JSON(500, gin.H{"error": "Delete operation failed"}) // DB error, return JSON error
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("User ID %d deleted successfully", id)}) // Success message JSON
}

// GetUsers returns the list of users (excluding password)
func GetUsers(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, username, email FROM users") // Fetch users
	if err != nil {
		c.JSON(500, gin.H{"error": "Data could not be retrieved"}) // If error, return JSON
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			c.JSON(500, gin.H{"error": "Row could not be read"}) // Row read error
			return
		}
		users = append(users, u)
	}

	c.JSON(200, users) // Return user list as JSON
}

// UpdateUser user update operation
func UpdateUser(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil { // Bind JSON body to User struct
		c.JSON(400, gin.H{"error": "Invalid data"}) // If invalid data, return JSON
		return
	}

	if u.Password != "" {
		_, err := Conn.Exec(context.Background(), "UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4", u.Username, u.Email, u.Password, u.ID) // Update including password
		if err != nil {
			c.JSON(500, gin.H{"error": "Update failed"})
			return
		}
	} else {
		_, err := Conn.Exec(context.Background(), "UPDATE users SET username=$1, email=$2 WHERE id=$3", u.Username, u.Email, u.ID) // Update without password
		if err != nil {
			c.JSON(500, gin.H{"error": "Update failed"})
			return
		}
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("User ID %d updated successfully", u.ID)}) // Success message
}

// CreateUser create new user
func CreateUser(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil { // Bind JSON to User struct
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	// Password hashing should be done here, currently storing as plain text

	_, err := Conn.Exec(context.Background(), "INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", u.Username, u.Password, u.Email) // Insert into DB
	if err != nil {
		c.JSON(500, gin.H{"error": "User could not be created"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"}) // Success message
}
