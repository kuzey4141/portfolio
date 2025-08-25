package routes

import (
	"portfolio/about"
	"portfolio/contact"
	"portfolio/home"
	"portfolio/middleware"
	"portfolio/projects"
	"portfolio/user"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *gin.Engine, dbPool *pgxpool.Pool) {
	// CORS settings - Fixed for credentials
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Always set specific origin, never wildcard with credentials
		if origin == "http://portfolio-kuzey-2025.s3-website.eu-central-1.amazonaws.com" ||
			origin == "https://portfolio-kuzey-2025.s3-website.eu-central-1.amazonaws.com" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			// For direct API calls without origin
			c.Header("Access-Control-Allow-Origin", "http://portfolio-kuzey-2025.s3-website.eu-central-1.amazonaws.com")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "false") // Set to false to avoid wildcard issue

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Set database connection for all packages
	home.SetDB(dbPool)
	about.SetDB(dbPool)
	projects.SetDB(dbPool)
	contact.SetDB(dbPool)
	user.SetDB(dbPool)

	// PUBLIC ROUTES
	publicAPI := r.Group("/api")
	{
		publicAPI.GET("/home", home.GetHomes)
		publicAPI.GET("/about", about.GetAbouts)
		publicAPI.GET("/projects", projects.GetProjects)
		publicAPI.POST("/contact", contact.CreateContact)
		publicAPI.POST("/login", user.Login)
	}

	// ADMIN ROUTES
	adminAPI := r.Group("/api/admin")
	adminAPI.Use(middleware.AuthMiddleware())
	{
		// Contact management
		adminAPI.GET("/contact", contact.GetContacts)
		adminAPI.DELETE("/contact/:id", contact.DeleteContact)
		adminAPI.PUT("/contact", contact.UpdateContact)

		// Home management
		adminAPI.GET("/home", home.GetHomes)
		adminAPI.POST("/home", home.CreateHome)
		adminAPI.PUT("/home", home.UpdateHome)
		adminAPI.DELETE("/home/:id", home.DeleteHome)

		// About management
		adminAPI.POST("/about", about.CreateAbout)
		adminAPI.PUT("/about", about.UpdateAbout)
		adminAPI.DELETE("/about/:id", about.DeleteAbout)

		// Project management
		adminAPI.POST("/projects", projects.CreateProject)
		adminAPI.PUT("/projects/:id", projects.UpdateProject)
		adminAPI.DELETE("/projects/:id", projects.DeleteProject)
		adminAPI.GET("/projects", projects.GetProjects) // Add this line!
	}

	// SUPER ADMIN ROUTES
	superAdminAPI := r.Group("/api/superadmin")
	superAdminAPI.Use(middleware.AuthMiddleware())
	superAdminAPI.Use(middleware.SuperAdminMiddleware())
	{
		superAdminAPI.GET("/users", user.GetUsers)
		superAdminAPI.POST("/users", user.CreateUser)
		superAdminAPI.PUT("/users", user.UpdateUser)
		superAdminAPI.DELETE("/users/:id", user.DeleteUser)
	}
}
