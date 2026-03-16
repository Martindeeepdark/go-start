package main

import (
	"context"
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
	"{{.Module}}/pkg/database"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	// 初始化数据库
	db, err := database.New(&cfg.Database)
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}
	defer db.Close()
	logger.Info("数据库连接成功")

	// 依赖注入
	repo := repository.New(db)
	svc := service.New(repo)
	ctrl := controller.New(svc)

	// 初始化路由
	r := gin.Default()
	ctrl.RegisterRoutes(r)

	// 启动服务
	addr := cfg.Server.Addr()
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logger.Info("服务启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务关闭异常", zap.Error(err))
	}
	logger.Info("服务已关闭")
}
