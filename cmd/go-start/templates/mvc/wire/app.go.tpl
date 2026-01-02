package wire

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ============================================
// App 结构体：应用入口
// ============================================

// App 应用结构体
//
// 说明：
// - 包含应用运行所需的所有依赖
// - Wire 会自动填充这些字段
// - 字段必须是导出的（大写开头）
type App struct {
	Config    *ConfigWrapper
	Logger    *zap.Logger
	DB        interface{}
	Cache     interface{}
	Router    *gin.Engine
}

// ConfigWrapper 配置包装器
// 使用 interface{} 避免循环依赖
type ConfigWrapper struct {
	Config interface{}
}

// ============================================
// App 方法：运行和关闭
// ============================================

// Run 运行应用
//
// 说明：
// - 启动 HTTP 服务器
// - 监听配置的端口
// - 阻塞直到发生错误
func (a *App) Run() error {
	// 获取配置（需要类型断言，因为使用 interface{}）
	// 实际项目中可以直接使用具体类型

	addr := ":8080" // 默认端口
	a.Logger.Info("Starting server",
		zap.String("addr", addr),
		zap.String("mode", gin.Mode()),
	)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	// 启动服务器（非阻塞）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	a.Logger.Info("Server started successfully")

	// 等待中断信号优雅关闭
	return a.waitForShutdown(srv)
}

// waitForShutdown 等待中断信号并优雅关闭
func (a *App) waitForShutdown(srv *http.Server) error {
	// 实现优雅关闭逻辑
	// 见下面的完整实现
	return nil
}

// Shutdown 优雅关闭应用
//
// 说明：
// - 关闭 HTTP 服务器
// - 关闭数据库连接
// - 关闭 Redis 连接
// - 刷新日志
func (a *App) Shutdown(ctx context.Context) error {
	a.Logger.Info("Shutting down server...")

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := srv.Shutdown(ctx); err != nil {
		a.Logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	// 关闭数据库
	if db, ok := a.DB.(interface{ Close() error }); ok {
		if err := db.Close(); err != nil {
			a.Logger.Error("Failed to close database", zap.Error(err))
		}
	}

	// 关闭 Redis
	if cache, ok := a.Cache.(interface{ Close() error }); ok {
		if err := cache.Close(); err != nil {
			a.Logger.Error("Failed to close cache", zap.Error(err))
		}
	}

	// 刷新日志
	a.Logger.Sync()

	a.Logger.Info("Server exited successfully")
	return nil
}
