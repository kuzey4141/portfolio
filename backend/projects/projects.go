package projects // projects package

import (
	"context" // context for DB operations
	"fmt"     // for printing to console
	"strconv" // for string-int conversions

	"github.com/gin-gonic/gin" // Gin framework

	"github.com/jackc/pgx/v5" // PostgreSQL library
)

// Project struct represents the data in the projects table
type Project struct {
	ID          int    `json:"id"`          // shown as id in JSON
	Name        string `json:"name"`        // shown as name in JSON
	Description string `json:"description"` // shown as description in JSON
	Message     string `json:"message"`     // shown as message in JSON
}

var Conn *pgx.Conn // Global database connection

func SetDB(conn *pgx.Conn) {
	Conn = conn // Set the global connection
}

// DeleteProject Gin handler: For delete operation
func DeleteProject(c *gin.Context) {
	idStr := c.Param("id")         // Get :id parameter from URL
	id, err := strconv.Atoi(idStr) // Convert string to int
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"}) // Return JSON error for invalid ID
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM projects WHERE id=$1", id) // Delete query
	if err != nil {
		c.JSON(500, gin.H{"error": "Delete operation failed"}) // Return JSON error if DB error
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Project ID %d deleted successfully", id)}) // Success message as JSON
}

// GetProjects returns all projects as JSON
func GetProjects(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, name, description, message FROM projects") // message field eklendi
	if err != nil {
		c.JSON(500, gin.H{"error": "Data could not be retrieved"}) // Return JSON error if failed
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Message); err != nil { // Message olarak değişti
			c.JSON(500, gin.H{"error": "Could not read row"}) // Row read error
			return
		}
		projects = append(projects, p)
	}

	c.JSON(200, projects) // Successfully return projects as JSON
}

// UpdateProject Gin handler for update operation
func UpdateProject(c *gin.Context) {
	idStr := c.Param("id")         // Get id from URL
	id, err := strconv.Atoi(idStr) // Convert string to int
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"}) // Return error if ID is invalid
		return
	}

	var p Project
	if err := c.BindJSON(&p); err != nil { // Bind JSON to Project struct
		c.JSON(400, gin.H{"error": "Invalid data"}) // JSON parse error
		return
	}

	sql := `UPDATE projects SET name=$1, description=$2, message=$3 WHERE id=$4`
	_, err = Conn.Exec(context.Background(), sql, p.Name, p.Description, p.Message, id) // Message olarak değişti
	if err != nil {
		c.JSON(500, gin.H{"error": "Update failed"}) // DB error
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Project ID %d updated successfully", id)}) // Success message
}

// CreateProject Gin handler for adding a new project
func CreateProject(c *gin.Context) {
	var p Project
	if err := c.BindJSON(&p); err != nil { // Bind JSON to Project struct
		c.JSON(400, gin.H{"error": "Invalid data"}) // JSON parse error
		return
	}

	_, err := Conn.Exec(context.Background(),
		"INSERT INTO projects (name, description, message) VALUES ($1, $2, $3)",
		p.Name, p.Description, p.Message) // Message olarak değişti
	if err != nil {
		c.JSON(500, gin.H{"error": "Record could not be added"}) // Return error if failed
		return
	}

	c.JSON(201, gin.H{"message": "Project record added successfully"}) // Success message
}
