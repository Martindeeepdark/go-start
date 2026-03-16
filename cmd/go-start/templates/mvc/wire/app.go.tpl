package wire

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"{{.Module}}/config"
)

// App 应用结构体，持有所有依赖
type App struct {
	Config *config.Config
	Logger *zap.Logger
	Router *gin.Engine
}

// Run 启动 HTTP 服务器，阻塞直到收到退出信号
func (a *App) Run() error {
	addr := a.Config.Server.Addr()
	a.Logger.Info("Starting server", zap.String("addr", addr))

	srv := &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-quit:
		a.Logger.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			a.Logger.Error("Server forced to shutdown", zap.Error(err))
			return err
		}
		a.Logger.Info("Server exited successfully")
		return nil
	}
}
