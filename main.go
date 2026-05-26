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
	r.POST("/api/speedtest", handlers.SubmitSpeedTest)
	
	// 测速端点
	r.GET("/api/speedtest/download", func(c *gin.Context) {
		sizeStr := c.DefaultQuery("size", "10") // 单位是 MB
		var sizeMB int
		s, err := strconv.Atoi(sizeStr)
		if err != nil || s <= 0 {
			sizeMB = 10
		} else {
			sizeMB = s
		}
		
		if sizeMB > 100 {
			sizeMB = 100
		}
		
		sizeBytes := sizeMB * 1024 * 1024
		
		// 生成随机数据
		data := make([]byte, sizeBytes)
		rand.Read(data)
		
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=\"test.bin\"")
		c.Data(http.StatusOK, "application/octet-stream", data)
	})
	
	r.POST("/api/speedtest/upload", func(c *gin.Context) {
		var totalSize int64 = 0
		buf := make([]byte, 1024*10)
		for {
			n, err := c.Request.Body.Read(buf)
			if n > 0 {
				totalSize += int64(n)
			}
			if err != nil {
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{"size": totalSize})
	})

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
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.String(http.StatusOK, string(data))
	})

	addr := ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
