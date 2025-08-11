package routes

import (
	"portfolio/about"
	"portfolio/contact"
	"portfolio/home"
	"portfolio/projects"
	"portfolio/user"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(r *gin.Engine, dbConn *pgx.Conn) {

	home.SetDB(dbConn)
	about.SetDB(dbConn)
	projects.SetDB(dbConn)
	contact.SetDB(dbConn)
	user.SetDB(dbConn)

	
	r.GET("/api/home", home.GetHomes)
	r.GET("/api/about", about.GetAbouts)
	r.GET("/api/projects", projects.GetProjects)
	r.GET("/api/contact", contact.GetContacts)
	r.GET("/api/user", user.GetUsers)

	
	r.DELETE("/api/home/:id", home.DeleteHome)
	r.DELETE("/api/about/:id", about.DeleteAbout)
	r.DELETE("/api/projects/:id", projects.DeleteProject)
	r.DELETE("/api/contact/:id", contact.DeleteContact)
	r.DELETE("/api/user/:id", user.DeleteUser)

	r.PUT("/api/home", home.UpdateHome)
	r.PUT("/api/about", about.UpdateAbout)
	r.PUT("/api/projects", projects.UpdateProject)
	r.PUT("/api/contact", contact.UpdateContact)
	r.PUT("/api/user", user.UpdateUser)

	r.POST("/api/home", home.CreateHome)
	r.POST("/api/about", about.CreateAbout)
	r.POST("/api/projects", projects.CreateProject)
	r.POST("/api/contact", contact.CreateContact)
	r.POST("/api/user", user.CreateUser)

}
