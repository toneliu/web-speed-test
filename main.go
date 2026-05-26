package main

import (
	"embed"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"speedtest/pkg/database"
	"speedtest/pkg/handlers"
	"speedtest/pkg/middleware"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/*
var frontendFS embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())
	
	if err := database.Init("speedtest.db"); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/api/login", handlers.Login)
	r.GET("/api/units", handlers.GetUnits)
	
	// 测速端点 - 参照 librespeed/speedtest-go 的实现
	r.GET("/api/speedtest/garbage", func(c *gin.Context) {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		
		sizeStr := c.DefaultQuery("ckSize", "4")
		size, _ := strconv.Atoi(sizeStr)
		if size <= 0 {
			size = 4
		}
		chunkSize := size * 1024 * 1024
		
		data := make([]byte, chunkSize)
		rand.Read(data)
		
		c.Data(http.StatusOK, "application/octet-stream", data)
	})
	
	r.GET("/api/speedtest/empty", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.JSON(http.StatusOK, gin.H{})
	})
	
	r.POST("/api/speedtest/empty", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.JSON(http.StatusOK, gin.H{})
	})
	
	r.POST("/api/speedtest/getIP", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.JSON(http.StatusOK, gin.H{
			"processedString": c.ClientIP(),
		})
	})
	
	r.GET("/api/speedtest/ping", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.JSON(http.StatusOK, gin.H{})
	})
	
	r.POST("/api/speedtest/results", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.JSON(http.StatusOK, gin.H{})
	})

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/speedtest", handlers.SubmitSpeedTest)
		auth.GET("/speedtests", handlers.GetSpeedTests)
		auth.GET("/stats", handlers.GetStats)
		auth.PUT("/user/password", handlers.ChangePassword)

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
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.String(http.StatusOK, string(data))
	})

	addr := ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
