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
	// PRIVATE ROUTES (Token gerektirir - Admin paneli)
	// =================
	adminAPI := r.Group("/api/admin")
	adminAPI.Use(middleware.AuthMiddleware()) // Bu middleware tüm admin route'larına uygulanır
	{
		// CONTACT - Gelen mesajları görme ve yönetme
		adminAPI.GET("/contact", contact.GetContacts)          // Gelen mesajları listele
		adminAPI.DELETE("/contact/:id", contact.DeleteContact) // Mesaj sil
		adminAPI.PUT("/contact", contact.UpdateContact)        // Mesaj güncelle (gerekirse)

		// HOME - Ana sayfa içeriği yönetimi
		adminAPI.GET("/home", home.GetHomes)          // Sadmin için de gerekebilir
		adminAPI.POST("/home", home.CreateHome)       // Yeni ana sayfa içeriği ekle
		adminAPI.PUT("/home", home.UpdateHome)        // Ana sayfa içeriğini güncelle
		adminAPI.DELETE("/home/:id", home.DeleteHome) // Ana sayfa içeriğini sil

		// ABOUT - Hakkımda bölümü yönetimi
		adminAPI.POST("/about", about.CreateAbout)       // Yeni hakkımda içeriği ekle
		adminAPI.PUT("/about", about.UpdateAbout)        // Hakkımda güncelle
		adminAPI.DELETE("/about/:id", about.DeleteAbout) // Hakkımda sil

		// PROJECTS - Proje yönetimi
		adminAPI.POST("/projects", projects.CreateProject)       // Yeni proje ekle
		adminAPI.PUT("/projects/:id", projects.UpdateProject)    // Proje güncelle
		adminAPI.DELETE("/projects/:id", projects.DeleteProject) // Proje sil

		// USERS - Kullanıcı yönetimi (Admin paneli)
		adminAPI.GET("/users", user.GetUsers)          // Kullanıcıları listele
		adminAPI.POST("/users", user.CreateUser)       // Yeni kullanıcı ekle
		adminAPI.PUT("/users", user.UpdateUser)        // Kullanıcı güncelle
		adminAPI.DELETE("/users/:id", user.DeleteUser) // Kullanıcı sil
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
