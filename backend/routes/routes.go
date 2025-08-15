package routes

import (
	"portfolio/about"
	"portfolio/contact"
	"portfolio/home"
	"portfolio/middleware"
	"portfolio/projects"
	"portfolio/user"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(r *gin.Engine, dbConn *pgx.Conn) {
	// Tüm package'lara database bağlantısını ver
	home.SetDB(dbConn)
	about.SetDB(dbConn)
	projects.SetDB(dbConn)
	contact.SetDB(dbConn)
	user.SetDB(dbConn)

	// =================
	// PUBLIC ROUTES (Token gerektirmez)
	// =================
	publicAPI := r.Group("/api")
	{
		// Portfolio bilgilerini herkes görebilir
		publicAPI.GET("/home", home.GetHomes)
		publicAPI.GET("/about", about.GetAbouts)
		publicAPI.GET("/projects", projects.GetProjects)

		// İletişim formu - sadece mesaj gönderme (herkes yapabilir)
		publicAPI.POST("/contact", contact.CreateContact)

		// Login endpoint - giriş yapmak için
		publicAPI.POST("/login", user.Login)
	}

	// =================
	// ADMIN ROUTES (Token gerektirir)
	// =================
	adminAPI := r.Group("/api/admin")
	adminAPI.Use(middleware.AuthMiddleware()) // Tüm admin route'larına auth middleware
	{
		// CONTACT - Gelen mesajları görme ve yönetme
		adminAPI.GET("/contact", contact.GetContacts)          // Gelen mesajları listele
		adminAPI.DELETE("/contact/:id", contact.DeleteContact) // Mesaj sil
		adminAPI.PUT("/contact", contact.UpdateContact)        // Mesaj güncelle

		// HOME - Ana sayfa içeriği yönetimi
		adminAPI.GET("/home", home.GetHomes)
		adminAPI.POST("/home", home.CreateHome)
		adminAPI.PUT("/home", home.UpdateHome)
		adminAPI.DELETE("/home/:id", home.DeleteHome)

		// ABOUT - Hakkımda bölümü yönetimi
		adminAPI.POST("/about", about.CreateAbout)
		adminAPI.PUT("/about", about.UpdateAbout)
		adminAPI.DELETE("/about/:id", about.DeleteAbout)

		// PROJECTS - Proje yönetimi
		adminAPI.POST("/projects", projects.CreateProject)
		adminAPI.PUT("/projects/:id", projects.UpdateProject)
		adminAPI.DELETE("/projects/:id", projects.DeleteProject)
	}

	// =================
	// SUPER ADMIN ROUTES (Sadece "admin" kullanıcısı)
	// =================
	superAdminAPI := r.Group("/api/superadmin")
	superAdminAPI.Use(middleware.AuthMiddleware())       // Önce auth kontrol
	superAdminAPI.Use(middleware.SuperAdminMiddleware()) // Sonra super admin kontrol
	{
		// USER MANAGEMENT - Sadece super admin yapabilir
		superAdminAPI.GET("/users", user.GetUsers)          // Kullanıcıları listele
		superAdminAPI.POST("/users", user.CreateUser)       // Yeni kullanıcı ekle
		superAdminAPI.PUT("/users", user.UpdateUser)        // Kullanıcı güncelle
		superAdminAPI.DELETE("/users/:id", user.DeleteUser) // Kullanıcı sil
	}

	// =================
	// CORS ayarları (Frontend için)
	// =================
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
