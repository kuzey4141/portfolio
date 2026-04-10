package server

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

// StartFrontend starts the frontend in development mode
func StartFrontend() {
	if gin.Mode() != gin.DebugMode {
		return
	}

	go func() {
		log.Println("Starting frontend...")
		cmd := exec.Command("bash", "-c", "cd ../frontend && npm start")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Frontend error: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)
}

// SetupStaticFiles configures static file serving in production mode
func SetupStaticFiles(r *gin.Engine) {
	if gin.Mode() != gin.ReleaseMode {
		return
	}

	r.Static("/static", "./frontend/build/static")
	r.StaticFile("/", "./frontend/build/index.html")
	r.StaticFile("/admin", "./frontend/build/index.html")
	r.StaticFile("/favicon.ico", "./frontend/build/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/build/index.html")
	})
}
