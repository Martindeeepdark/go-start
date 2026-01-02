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
	"{{.Module}}/config"
	"{{.Module}}/internal/controller"
	"{{.Module}}/internal/repository"
	"{{.Module}}/internal/service"
	"{{.Module}}/pkg/cache"
	"{{.Module}}/pkg/database"
	httpx "{{.Module}}/pkg/httpx/middleware"
	"{{.Module}}/pkg/httpx/response"
	"{{.Module}}/pkg/httpx/router"
	"go.uber.org/zap"
	{{if .WithSwagger}}
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	{{end}}
)

// @title           {{.ProjectName}} API
// @version         1.0
// @description     {{if .Description}}{{.Description}}{{else}}{{.ProjectName}} 服务的 API 文档{{end}}
// @host            localhost:{{.ServerPort}}
// @BasePath        /api/v1

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
	db, err := database.New(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Database connected successfully")

	{{if .WithRedis}}
	// Initialize Redis
	cacheClient, err := cache.New(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect redis", zap.Error(err))
	}
	defer cacheClient.Close()
	logger.Info("Redis connected successfully")
	{{else}}
	var cacheClient *cache.Cache
	{{end}}

	// ============================================
	// 依赖注入链 (Dependency Injection)
	// ============================================
	//
	// 说明：
	//   1. 我们使用简单的手动依赖注入，不使用 Wire/Fx 等库
	//   2. 依赖关系清晰：DB → Repository → Service → Controller
	//   3. 构造顺序：从底层到顶层
	//   4. 使用接口而不是具体类型，便于测试和切换实现
	//
	// 依赖图：
	//   ┌─────────────┐
	//   │ Controller  │  ← HTTP 请求处理层
	//   └──────┬──────┘
	//          │ 依赖
	//   ┌──────▼──────┐
	//   │   Service   │  ← 业务逻辑层
	//   └──────┬──────┘
	//          │ 依赖
	//   ┌──────▼──────┐
	//   │ Repository  │  ← 数据访问层
	//   └──────┬──────┘
	//          │ 依赖
	//   ┌──────▼──────┐
	//   │     DB      │  ← 数据库
	//   └─────────────┘
	//
	// 为什么使用接口？
	//   ✅ 易于测试：可以 mock Repository
	//   ✅ 灵活切换：可以替换不同实现
	//   ✅ 解耦合：Service 不依赖具体实现
	//
	// 示例：测试时 mock Repository
	//   mockRepo := &MockUserRepository{}
	//   service := service.NewUserService(mockRepo, nil)
	//
	// 示例：切换到不同的数据存储
	//   var repo repository.UserRepository
	//   if useRedis {
	//       repo = repository.NewRedisUserRepository(redis)
	//   } else {
	//       repo = repository.NewMySQLUserRepository(db)
	//   }

	// Step 1: 初始化 Repository 层（依赖 DB）
	// 说明：Repository 返回接口类型，不是具体类型
	repo := repository.New(db)

	// Step 2: 初始化 Service 层（依赖 Repository 和 Cache）
	// 说明：Service 依赖 Repository 接口，符合依赖倒置原则（DIP）
	svc := service.New(repo, cacheClient)

	// Step 3: 初始化 Controller 层（依赖 Service）
	// 说明：Controller 依赖 Service，处理 HTTP 请求
	ctrl := controller.New(svc)

	// ============================================
	// 依赖注入完成
	// ============================================

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

	srv := &http.Server{
		Addr:    addr,
		Handler: r.Engine(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited successfully")
}

func registerRoutes(r *router.Router, ctrl *controller.Controller) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	{{if .WithSwagger}}
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	{{end}}

	// API v1
	_ = r.Group("/api/v1")
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
