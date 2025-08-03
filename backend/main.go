package main

import (
    "portfolio/about"
    "portfolio/contact"
    "portfolio/db"
    "portfolio/home"
    "portfolio/projects"
    "portfolio/user"

    "github.com/gin-gonic/gin"
)

func main() {
    db.ConnectDB()
    defer db.Conn.Close(context.Background())

    home.SetDB(db.Conn)
    about.SetDB(db.Conn)
    projects.SetDB(db.Conn)
    contact.SetDB(db.Conn)
    user.SetDB(db.Conn)

    r := gin.Default()

    // GET Endpoints
    r.GET("/api/home", home.GetHomes)
    r.GET("/api/about", about.GetAbouts)
    r.GET("/api/projects", projects.GetProjects)
    r.GET("/api/contact", contact.GetContacts)
    r.GET("/api/user", user.GetUsers)

    // DELETE Endpoints
    r.DELETE("/api/home/:id", home.DeleteHome)
    r.DELETE("/api/about/:id", about.DeleteAbout)
    r.DELETE("/api/projects/:id", projects.DeleteProject)
    r.DELETE("/api/contact/:id", contact.DeleteContact)
    r.DELETE("/api/user/:id", user.DeleteUser)

    // PUT Endpoints
    r.PUT("/api/home", home.UpdateHome)
    r.PUT("/api/about", about.UpdateAbout)
    r.PUT("/api/projects", projects.UpdateProject)
    r.PUT("/api/contact", contact.UpdateContact)
    r.PUT("/api/user", user.UpdateUser)

    // POST Endpoints
    r.POST("/api/home", home.CreateHome)
    r.POST("/api/about", about.CreateAbout)
    r.POST("/api/projects", projects.CreateProject)
    r.POST("/api/contact", contact.CreateContact)
    r.POST("/api/user", user.CreateUser)

    r.Run(":8080") // 8080 portunda çalıştır
}
