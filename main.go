package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
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

	gin.SetMode(gin.ReleaseMode)
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

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		var data []byte
		var contentType string

		if path == "/" || !strings.Contains(path, ".") {
			data, err = fs.ReadFile(staticFS, "index.html")
			contentType = "text/html; charset=utf-8"
		} else {
			data, err = fs.ReadFile(staticFS, strings.TrimPrefix(path, "/"))
			if err != nil {
				data, err = fs.ReadFile(staticFS, "index.html")
				contentType = "text/html; charset=utf-8"
			} else {
				contentType = http.DetectContentType(data)
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", contentType)
		c.String(http.StatusOK, string(data))
	})

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
