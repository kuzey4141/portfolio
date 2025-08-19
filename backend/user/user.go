package user

import (
	"context"
	"fmt"
	"net/http"
	"portfolio/auth" // Auth package import
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL pool library
)

// User struct represents data in the users table
type User struct {
	ID       int    `json:"id"`                 // ID field in JSON
	Username string `json:"username"`           // Username field in JSON
	Password string `json:"password,omitempty"` // Password in JSON (shown empty if not provided)
	Email    string `json:"email"`              // Email in JSON
}

// LoginRequest struct for login endpoint
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse struct for login response
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    User   `json:"user"`
}

var Pool *pgxpool.Pool

func SetDB(pool *pgxpool.Pool) {
	Pool = pool
}

// Login function - Authentication login function
func Login(c *gin.Context) {
	var loginReq LoginRequest

	// Get data from JSON
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	// Find user in database
	var user User
	err := Pool.QueryRow(context.Background(),
		"SELECT id, username, password, email FROM users WHERE username=$1",
		loginReq.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)

	if err != nil {
		fmt.Println("User not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check password
	if err := auth.CheckPassword(user.Password, loginReq.Password); err != nil {
		fmt.Println("Wrong password:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		fmt.Println("Token generation error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Don't send password in response
	user.Password = ""

	// Successful login response
	response := LoginResponse{
		Message: "Login successful",
		Token:   token,
		User:    user,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser delete operation
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = Pool.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Delete operation failed"})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("User ID %d deleted successfully", id)})
}

// GetUsers returns the list of users (excluding password)
func GetUsers(c *gin.Context) {
	rows, err := Pool.Query(context.Background(), "SELECT id, username, email FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": "Data could not be retrieved"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			c.JSON(500, gin.H{"error": "Row could not be read"})
			return
		}
		users = append(users, u)
	}

	c.JSON(200, users)
}

// UpdateUser user update operation
func UpdateUser(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	// If password exists, hash it
	if u.Password != "" {
		hashedPassword, err := auth.HashPassword(u.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not hash password"})
			return
		}

		_, err = Pool.Exec(context.Background(),
			"UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4",
			u.Username, u.Email, hashedPassword, u.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Update failed"})
			return
		}
	} else {
		_, err := Pool.Exec(context.Background(),
			"UPDATE users SET username=$1, email=$2 WHERE id=$3",
			u.Username, u.Email, u.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Update failed"})
			return
		}
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("User ID %d updated successfully", u.ID)})
}

// CreateUser create new user
func CreateUser(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not hash password"})
		return
	}

	_, err = Pool.Exec(context.Background(),
		"INSERT INTO users (username, password, email) VALUES ($1, $2, $3)",
		u.Username, hashedPassword, u.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "User could not be created"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}
