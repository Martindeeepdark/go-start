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
	"github.com/github.com/yourname/test-api/config"
	"github.com/github.com/yourname/test-api/internal/controller"
	"github.com/github.com/yourname/test-api/internal/repository"
	"github.com/github.com/yourname/test-api/internal/service"
	"github.com/github.com/yourname/test-api/pkg/cache"
	"github.com/github.com/yourname/test-api/pkg/database"
	httpx "github.com/github.com/yourname/test-api/pkg/httpx/middleware"
	"github.com/github.com/yourname/test-api/pkg/httpx/response"
	"github.com/github.com/yourname/test-api/pkg/httpx/router"
	"go.uber.org/zap"
)

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

	// Initialize repository layer
	repo := repository.New(db)

	// Initialize service layer
	svc := service.New(repo, cacheClient)

	// Initialize controller layer
	ctrl := controller.New(svc)

	// Initialize router
	r := router.New()
	r.Use(
		httpx.Logger(logger),
		httpx.Recovery(logger),
		httpx.CORS(),
		httpx.RequestID(),
	)

	// Register routes
	registerRoutes(r, ctrl)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Starting server", zap.String("addr", addr))

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

func registerRoutes(r *router.Router, ctrl *controller.Controller) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Example: User routes
		// user := v1.Group("/users")
		// {
		// 	user.GET("", ctrl.User.List)
		// 	user.GET("/:id", ctrl.User.Get)
		// 	user.POST("", ctrl.User.Create)
		// 	user.PUT("/:id", ctrl.User.Update)
		// 	user.DELETE("/:id", ctrl.User.Delete)
		// }
	}
}
