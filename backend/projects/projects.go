package projects

import (
	"context" // Context for DB operations
	"fmt"     // For printing to console
	"strconv" // For string-int conversions

	"github.com/gin-gonic/gin"        // Gin framework
	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL pool library
)

// Project struct represents the data in the projects table
type Project struct {
	ID           int    `json:"id"`           // Shown as id in JSON
	Name         string `json:"name"`         // Shown as name in JSON
	Description  string `json:"description"`  // Shown as description in JSON
	Message      string `json:"message"`      // Shown as message in JSON
	ImageURL     string `json:"image_url"`    // New field for project image
	Technologies string `json:"technologies"` // New field for tech stack
	GithubURL    string `json:"github_url"`   // New field for GitHub link
	DemoURL      string `json:"demo_url"`     // New field for demo link
}

var Pool *pgxpool.Pool // Global database pool

func SetDB(pool *pgxpool.Pool) {
	Pool = pool // Set the global pool
}

// DeleteProject Gin handler: For delete operation
func DeleteProject(c *gin.Context) {
	idStr := c.Param("id")         // Get :id parameter from URL
	id, err := strconv.Atoi(idStr) // Convert string to int
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"}) // Return JSON error for invalid ID
		return
	}

	_, err = Pool.Exec(context.Background(), "DELETE FROM projects WHERE id=$1", id) // Delete query
	if err != nil {
		fmt.Println("Delete error:", err)
		c.JSON(500, gin.H{"error": "Delete operation failed"}) // Return JSON error if DB error
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Project ID %d deleted successfully", id)}) // Success message as JSON
}

// GetProjects returns all projects as JSON
func GetProjects(c *gin.Context) {
	rows, err := Pool.Query(context.Background(),
		"SELECT id, name, description, message, COALESCE(image_url, ''), COALESCE(technologies, ''), COALESCE(github_url, ''), COALESCE(demo_url, '') FROM projects ORDER BY id DESC") // Updated query with new fields
	if err != nil {
		fmt.Println("Query error:", err)
		c.JSON(500, gin.H{"error": "Data could not be retrieved"}) // Return JSON error if failed
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Message, &p.ImageURL, &p.Technologies, &p.GithubURL, &p.DemoURL); err != nil { // Updated scan with new fields
			fmt.Println("Row scan error:", err)
			c.JSON(500, gin.H{"error": "Could not read row"}) // Row read error
			return
		}
		projects = append(projects, p)
	}

	fmt.Printf("Found %d projects\n", len(projects)) // Debug log
	c.JSON(200, projects)                            // Successfully return projects as JSON
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
	if err := c.ShouldBindJSON(&p); err != nil { // Bind JSON to Project struct
		fmt.Println("JSON bind error:", err)
		c.JSON(400, gin.H{"error": "Invalid data"}) // JSON parse error
		return
	}

	sql := `UPDATE projects SET name=$1, description=$2, message=$3, image_url=$4, technologies=$5, github_url=$6, demo_url=$7 WHERE id=$8`
	_, err = Pool.Exec(context.Background(), sql, p.Name, p.Description, p.Message, p.ImageURL, p.Technologies, p.GithubURL, p.DemoURL, id) // Updated query with new fields
	if err != nil {
		fmt.Println("Update error:", err)
		c.JSON(500, gin.H{"error": "Update failed"}) // DB error
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Project ID %d updated successfully", id)}) // Success message
}

// CreateProject Gin handler for adding a new project
func CreateProject(c *gin.Context) {
	var p Project
	if err := c.ShouldBindJSON(&p); err != nil { // Bind JSON to Project struct
		fmt.Println("JSON bind error:", err)
		c.JSON(400, gin.H{"error": "Invalid data"}) // JSON parse error
		return
	}

	_, err := Pool.Exec(context.Background(),
		"INSERT INTO projects (name, description, message, image_url, technologies, github_url, demo_url) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		p.Name, p.Description, p.Message, p.ImageURL, p.Technologies, p.GithubURL, p.DemoURL) // Updated insert with new fields
	if err != nil {
		fmt.Println("Insert error:", err)
		c.JSON(500, gin.H{"error": "Record could not be added"}) // Return error if failed
		return
	}

	c.JSON(201, gin.H{"message": "Project record added successfully"}) // Success message
}
