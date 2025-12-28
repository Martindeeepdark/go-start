package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/{{.Module}}/config"
	"github.com/{{.Module}}/internal/application"
	"github.com/{{.Module}}/internal/infrastructure/persistence"
	"github.com/{{.Module}}/internal/interface/http"
	"github.com/{{.Module}}/pkg/cache"
	"github.com/{{.Module}}/pkg/database"
	httpx "github.com/{{.Module}}/pkg/httpx/middleware"
	"github.com/{{.Module}}/pkg/httpx/response"
	"github.com/{{.Module}}/pkg/httpx/router"
	"go.uber.org/zap"
)

// @title           {{.ProjectName}} API (DDD)
// @version         1.0
// @description     {{if .Description}}{{.Description}}{{else}}{{.ProjectName}} 服务的 API 文档（DDD 架构）{{end}}
// @host            localhost:{{.ServerPort}}
// @BasePath        /api

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Database connected successfully")

	// Initialize Redis
	cacheClient, err := cache.New(cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect redis", zap.Error(err))
	}
	defer cacheClient.Close()
	logger.Info("Redis connected successfully")

	// Initialize Infrastructure Layer (Repositories)
	repos := persistence.NewRepositories(db, cacheClient)
	logger.Info("Infrastructure layer initialized")

	// Initialize Application Layer (Services)
	applications := application.NewApplications(repos)
	logger.Info("Application layer initialized")

	// Initialize Interface Layer (Controllers)
	controllers := http.NewControllers(applications)
	logger.Info("Interface layer initialized")

	// Initialize router
	r := router.New()
	r.Use(
		httpx.Logger(logger),
		httpx.Recovery(logger),
		httpx.CORS(),
		httpx.RequestID(),
	)

	// Register routes
	http.RegisterRoutes(r, controllers)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok", "architecture": "DDD"})
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Starting DDD server", zap.String("addr", addr))

	go func() {
		if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.Engine().Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited successfully")
}
