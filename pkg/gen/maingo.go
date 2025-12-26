package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateMainGo 生成可运行的 main.go
func (g *DatabaseGenerator) GenerateMainGo() error {
	modulePath := getModulePath(g.config.Module)
	tables := g.config.Tables

	// 创建 cmd/server 目录
	outputDir := filepath.Join(g.config.Output, "cmd", "server")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建 cmd/server 目录失败: %w", err)
	}

	outputPath := filepath.Join(outputDir, "main.go")

	// 检查文件是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Println("⚠️  cmd/server/main.go 已存在，跳过生成")
		return nil
	}

	fmt.Println("📦 正在生成 cmd/server/main.go...")

	// 准备模型名称列表
	var modelNames []string
	for _, table := range tables {
		modelNames = append(modelNames, toModelName(table))
	}

	// 渲染模板
	if err := g.renderMainGoTemplate(outputPath, modulePath, modelNames); err != nil {
		return err
	}

	fmt.Println("     ✓ cmd/server/main.go 创建成功")

	// 生成配置文件
	if err := g.GenerateConfigYAML(); err != nil {
		return err
	}

	return nil
}

// renderMainGoTemplate 渲染 main.go 模板
func (g *DatabaseGenerator) renderMainGoTemplate(outputPath, modulePath string, modelNames []string) error {
	tmpl := `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"{{.ModulePath}}/internal/controller"
	"{{.ModulePath}}/internal/repository"
	"{{.ModulePath}}/internal/routes"
	"{{.ModulePath}}/internal/service"
	"{{.ModulePath}}/pkg/cache"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting {{.ModulePath}}...")

	// 从环境变量读取数据库配置
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		// 默认值
		dsn = "root:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// 初始化数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		zapLogger.Fatal("Failed to connect database", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	defer sqlDB.Close()
	zapLogger.Info("Database connected successfully")

	// 初始化 Redis (可选)
	cacheClient := cache.New()

	// 初始化 Repository 层
	{{range .ModelInfos}}
	{{.LowerCamelCase}}Repo := repository.New{{.Name}}Repository(db)
	{{end}}

	// 初始化 Service 层
	{{range .ModelInfos}}
	{{.LowerCamelCase}}Service := service.New{{.Name}}Service({{.LowerCamelCase}}Repo, cacheClient)
	{{end}}

	// 初始化 Controller 层
	controllers := &routes.Controllers{
		{{range .ModelInfos}}
		{{.Name}}: controller.New{{.Name}}Controller({{.LowerCamelCase}}Service),
		{{end}}
	}

	// 初始化路由
	r := gin.Default()

	// 注册自动生成的路由
	routes.RegisterAutoRoutes(r, controllers)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 启动 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		zapLogger.Info("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server exited")
}
`

	t, err := template.New("main").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	// 准备模型名称和小驼峰名称的映射
	type ModelInfo struct {
		Name            string
		LowerCamelCase  string
	}
	var modelInfos []ModelInfo
	for _, name := range modelNames {
		modelInfos = append(modelInfos, ModelInfo{
			Name:           name,
			LowerCamelCase: toLowerCamelCaseMain(name),
		})
	}

	data := map[string]interface{}{
		"ModulePath":  modulePath,
		"ModelNames":  modelNames,
		"ModelInfos":  modelInfos,
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	return nil
}

// GenerateConfigYAML 生成配置文件示例
func (g *DatabaseGenerator) GenerateConfigYAML() error {
	outputPath := filepath.Join(g.config.Output, "config.yaml.example")

	// 检查文件是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Println("⚠️  config.yaml.example 已存在，跳过生成")
		return nil
	}

	fmt.Println("📦 正在生成 config.yaml.example...")

	content := `# 数据库配置
DATABASE_DSN: "root:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

# 服务器配置
SERVER_PORT: "8080"
SERVER_MODE: "release" # debug, release

# Redis 配置 (可选)
REDIS_ADDR: "localhost:6379"
REDIS_PASSWORD: ""
REDIS_DB: "0"

# 日志配置
LOG_LEVEL: "info" # debug, info, warn, error
`

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 config.yaml.example 失败: %w", err)
	}

	fmt.Println("     ✓ config.yaml.example 创建成功")
	return nil
}

// toLowerCamelCaseMain 转换为小驼峰命名（用于 main.go 生成）
func toLowerCamelCaseMain(s string) string {
	if len(s) == 0 {
		return ""
	}
	// Users -> users
	// UserProfile -> userProfile
	return strings.ToLower(s[:1]) + s[1:]
}
