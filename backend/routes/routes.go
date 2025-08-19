package routes

import (
	"portfolio/about"
	"portfolio/contact"
	"portfolio/home"
	"portfolio/middleware"
	"portfolio/projects"
	"portfolio/user"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool" // Import pgxpool
)

func SetupRoutes(r *gin.Engine, dbPool *pgxpool.Pool) { // Changed parameter type
	// CORS settings - MUST BE AT THE TOP
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "false")

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Set database connection for all packages
	home.SetDB(dbPool)     // Pass pool instead of conn
	about.SetDB(dbPool)    // Pass pool instead of conn
	projects.SetDB(dbPool) // Pass pool instead of conn
	contact.SetDB(dbPool)  // Pass pool instead of conn
	user.SetDB(dbPool)     // Pass pool instead of conn

	// =================
	// PUBLIC ROUTES (No token required)
	// =================
	publicAPI := r.Group("/api")
	{
		// Portfolio information accessible to everyone
		publicAPI.GET("/home", home.GetHomes)
		publicAPI.GET("/about", about.GetAbouts)
		publicAPI.GET("/projects", projects.GetProjects)

		// Contact form - message sending (accessible to everyone)
		publicAPI.POST("/contact", contact.CreateContact)

		// Login endpoint
		publicAPI.POST("/login", user.Login)
	}

	// =================
	// ADMIN ROUTES (Token required)
	// =================
	adminAPI := r.Group("/api/admin")
	adminAPI.Use(middleware.AuthMiddleware()) // Auth middleware for all admin routes
	{
		// CONTACT - View and manage incoming messages
		adminAPI.GET("/contact", contact.GetContacts)          // List incoming messages
		adminAPI.DELETE("/contact/:id", contact.DeleteContact) // Delete message
		adminAPI.PUT("/contact", contact.UpdateContact)        // Update message

		// HOME - Home page content management
		adminAPI.GET("/home", home.GetHomes)
		adminAPI.POST("/home", home.CreateHome)
		adminAPI.PUT("/home", home.UpdateHome)
		adminAPI.DELETE("/home/:id", home.DeleteHome)

		// ABOUT - About section management
		adminAPI.POST("/about", about.CreateAbout)
		adminAPI.PUT("/about", about.UpdateAbout)
		adminAPI.DELETE("/about/:id", about.DeleteAbout)

		// PROJECTS - Project management
		adminAPI.POST("/projects", projects.CreateProject)
		adminAPI.PUT("/projects/:id", projects.UpdateProject)
		adminAPI.DELETE("/projects/:id", projects.DeleteProject)
	}

	// =================
	// SUPER ADMIN ROUTES (Only "admin" user)
	// =================
	superAdminAPI := r.Group("/api/superadmin")
	superAdminAPI.Use(middleware.AuthMiddleware())       // First auth check
	superAdminAPI.Use(middleware.SuperAdminMiddleware()) // Then super admin check
	{
		// USER MANAGEMENT - Only super admin can do this
		superAdminAPI.GET("/users", user.GetUsers)          // List users
		superAdminAPI.POST("/users", user.CreateUser)       // Add new user
		superAdminAPI.PUT("/users", user.UpdateUser)        // Update user
		superAdminAPI.DELETE("/users/:id", user.DeleteUser) // Delete user
	}
}
