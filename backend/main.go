package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"xinhang-backend/cache"
	"xinhang-backend/config"
	"xinhang-backend/database"
	"xinhang-backend/email"
	"xinhang-backend/handlers"
	"xinhang-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	database.Connect(cfg)
	cache.ConnectRedis(cfg)
	email.Init(cfg)

	handlers.SetJWTSecret(cfg.JWTSecret)
	handlers.InitQR()

	r := gin.Default()
	r.Use(middleware.IPProtection())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(maxBodySize(1 << 20)) // 1MB

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/send-code",
			middleware.RateLimit(5, time.Minute),
			handlers.SendVerificationCode,
		)
		api.POST("/register",
			middleware.RateLimit(10, time.Minute),
			handlers.Register,
		)
		api.POST("/login",
			middleware.RateLimit(20, time.Minute),
			handlers.Login,
		)
		api.POST("/apply",
			middleware.RateLimit(5, time.Minute),
			optionalAuth(cfg.JWTSecret),
			handlers.Apply,
		)
		api.GET("/applications",
			middleware.JWTAuth(cfg.JWTSecret),
			middleware.AdminOnly(),
			handlers.ListApplications,
		)
		api.GET("/permit-qr", handlers.GeneratePermitQR)
		api.POST("/verify-qr", handlers.VerifyPermitQR)
		api.GET("/verify-qr", handlers.VerifyPermitQRGet)
		api.GET("/query-permit", handlers.QueryPermit)
	}

	if _, err := os.Stat("./dist"); err == nil {
		r.Static("/assets", "./dist/assets")
		r.Static("/images", "./dist/images")
		r.StaticFile("/favicon.ico", "./dist/favicon.ico")
		r.NoRoute(func(c *gin.Context) {
			c.File("./dist/index.html")
		})
		log.Println("Serving frontend from ./dist")
	}

	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func optionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Next()
			return
		}
		middleware.JWTAuth(secret)(c)
	}
}

func maxBodySize(n int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, n)
		c.Next()
	}
}
