package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"speedtest/pkg/database"
	"speedtest/pkg/handlers"
	"speedtest/pkg/middleware"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/*
var frontendFS embed.FS

func main() {
	if err := database.Init("speedtest.db"); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	r := gin.Default()

	r.POST("/api/login", handlers.Login)
	r.GET("/api/units", handlers.GetUnits)
	r.POST("/api/speedtest", handlers.SubmitSpeedTest)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/speedtests", handlers.GetSpeedTests)
		auth.GET("/stats", handlers.GetStats)

		admin := auth.Group("")
		admin.Use(middleware.AdminRequired())
		{
			admin.POST("/units", handlers.CreateUnit)
			admin.PUT("/units/:id", handlers.UpdateUnit)
			admin.DELETE("/units/:id", handlers.DeleteUnit)
			admin.GET("/units/:id/users", handlers.GetUnitUsers)
			admin.POST("/units/:id/users", handlers.CreateUnitUser)
			admin.PUT("/users/:user_id/password", handlers.ResetUserPassword)
			admin.DELETE("/users/:user_id", handlers.DeleteUser)

			admin.GET("/topology", handlers.GetTopologyLinks)
			admin.POST("/topology", handlers.CreateTopologyLink)
			admin.DELETE("/topology/:id", handlers.DeleteTopologyLink)
		}
	}

	staticFS, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		log.Fatal(err)
	}
	staticFileServer := http.FileServer(http.FS(staticFS))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/" || path == "/index.html" {
			c.FileFromFS("/index.html", http.FS(staticFS))
		} else {
			_, err := staticFS.Open(path)
			if err == nil {
				staticFileServer.ServeHTTP(c.Writer, c.Request)
			} else {
				c.FileFromFS("/index.html", http.FS(staticFS))
			}
		}
	})

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
