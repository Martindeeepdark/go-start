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

	// 生成 application 包
	if err := g.GenerateApplicationPackage(modulePath, modelNames); err != nil {
		return err
	}

	// 生成配置文件
	if err := g.GenerateConfigYAML(); err != nil {
		return err
	}

	// 🔥 新增：生成完整项目文件
	if err := g.generateProjectFiles(modulePath, modelNames); err != nil {
		return err
	}

	return nil
}

// renderMainGoTemplate 渲染 main.go 模板
func (g *DatabaseGenerator) renderMainGoTemplate(outputPath, modulePath string, modelNames []string) error {
	tmpl := `package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"{{.ModulePath}}/internal/application"
	"{{.ModulePath}}/internal/routes"
	"github.com/gin-gonic/gin"
)

// @title           {{.ModuleName}} API
// @version         1.0
// @description     {{.Description}}
func main() {
	ctx := context.Background()

	// Please do not change the order of the function calls below
	setCrashOutput()

	if err := loadEnv(); err != nil {
		panic("loadEnv failed, err=" + err.Error())
	}

	setLogLevel()

	if err := application.Init(ctx); err != nil {
		panic("InitializeInfra failed, err=" + err.Error())
	}

	startHttpServer()
}

// setCrashOutput 设置崩溃输出
func setCrashOutput() {
	crashFile, err := os.Create("crash.log")
	if err != nil {
		log.Printf("创建崩溃日志文件失败: %v", err)
		return
	}
	debug.SetCrashOutput(crashFile, debug.CrashOptions{})
}

// loadEnv 加载环境变量
func loadEnv() error {
	appEnv := os.Getenv("APP_ENV")
	fileName := ".env"
	if appEnv != "" {
		fileName = ".env." + appEnv
	}

	log.Printf("加载环境变量文件: %s", fileName)

	// 这里使用简单的环境变量读取
	// 实际项目中可以使用 godotenv.Load(fileName)
	// 但为了减少依赖，我们直接依赖 os.Getenv

	return nil
}

// setLogLevel 设置日志级别
func setLogLevel() {
	level := getEnv("LOG_LEVEL", "info")
	log.Printf("日志级别: %s", level)

	// 根据环境变量设置日志级别
	// 实际项目中可以集成 zap 等日志库
	_ = level // 占位符
}

// getEnv 获取环境变量，支持默认值
func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

// startHttpServer 启动 HTTP 服务器
func startHttpServer() {
	addr := getEnv("SERVER_ADDR", ":8080")
	readTimeout := getDurationEnv("READ_TIMEOUT", 15*time.Second)
	writeTimeout := getDurationEnv("WRITE_TIMEOUT", 15*time.Second)
	idleTimeout := getDurationEnv("IDLE_TIMEOUT", 60*time.Second)

	log.Printf("启动 HTTP 服务器: %s", addr)

	// 创建 Gin engine
	r := gin.Default()

	// 注册路由
	routes.RegisterRoutes(r)

	// TODO: 配置服务器超时、TLS 等
	// 这里简化处理，实际项目中可以使用类似 server.Start(r) 的方式

	// 启动服务器
	if err := r.Run(addr); err != nil {
		panic("启动 HTTP 服务器失败: " + err.Error())
	}
}

// getDurationEnv 获取时间间隔类型的环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	return defaultValue
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

// generateProjectFiles 生成完整项目配套文件
func (g *DatabaseGenerator) generateProjectFiles(modulePath string, modelNames []string) error {
	fmt.Println("\n📦 正在生成项目配套文件...")

	// 1. README.md
	if err := g.generateReadme(modulePath, modelNames); err != nil {
		return err
	}

	// 2. Makefile
	if err := g.generateMakefile(); err != nil {
		return err
	}

	// 3. .gitignore
	if err := g.generateGitignore(); err != nil {
		return err
	}

	// 4. .env.example
	if err := g.generateEnvExample(); err != nil {
		return err
	}

	// 5. scripts/test.sh
	if err := g.generateTestScript(modelNames); err != nil {
		return err
	}

	// 6. docker-compose.yml
	if err := g.generateDockerCompose(); err != nil {
		return err
	}

	fmt.Println("     ✓ 所有项目文件生成完成")
	return nil
}

// generateReadme 生成 README.md
func (g *DatabaseGenerator) generateReadme(modulePath string, modelNames []string) error {
	outputPath := filepath.Join(g.config.Output, "README.md")

	// 检查是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	projectName := filepath.Base(g.config.Output)
	apiEndpoints := g.generateAPIEndpoints(modelNames)

	readmeContent := "# " + projectName + `

这是使用 [go-start](https://github.com/Martindeeepdark/go-start) 生成的 Go API 项目。

## 快速开始

` + "```bash" + `
# 1. 复制配置文件
cp config.yaml.example config.yaml

# 2. 安装依赖
go mod download

# 3. 运行
make run
` + "```" + `

## API 端点

` + apiEndpoints + `

## 开发

` + "```bash" + `
make test     # 运行测试
make build    # 编译
make clean    # 清理
` + "```" + `

## 技术栈

- Go 1.21+
- Gin
- GORM Gen
- MySQL/PostgreSQL
- Redis (可选)

---
Generated by [go-start](https://github.com/Martindeeepdark/go-start)
`

	return os.WriteFile(outputPath, []byte(readmeContent), 0644)
}

// generateAPIEndpoints 生成 API 端点说明
func (g *DatabaseGenerator) generateAPIEndpoints(modelNames []string) string {
	var endpoints string
	endpoints += "- 健康检查: `GET /health`\n\n"

	for _, name := range modelNames {
		lowerName := toLowerCamelCaseMain(name)
		endpoints += fmt.Sprintf("### %s\n", name)
		endpoints += fmt.Sprintf("- 获取列表: `GET /api/v1/%s`\n", lowerName)
		endpoints += fmt.Sprintf("- 获取详情: `GET /api/v1/%s/:id`\n", lowerName)
		endpoints += fmt.Sprintf("- 创建: `POST /api/v1/%s`\n", lowerName)
		endpoints += fmt.Sprintf("- 更新: `PUT /api/v1/%s/:id`\n", lowerName)
		endpoints += fmt.Sprintf("- 删除: `DELETE /api/v1/%s/:id`\n", lowerName)
		endpoints += "\n"
	}

	return endpoints
}

// generateMakefile 生成 Makefile
func (g *DatabaseGenerator) generateMakefile() error {
	outputPath := filepath.Join(g.config.Output, "Makefile")

	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	content := `.PHONY: run build test clean mod-tidy help

APP_NAME := $(shell basename $(PWD))
GO := go

run:
	@echo "🚀 启动服务..."
	@$(GO) run cmd/server/main.go

build:
	@echo "🔨 编译..."
	@mkdir -p bin
	@$(GO) build -o bin/$(APP_NAME) cmd/server/main.go

test:
	@echo "🧪 运行测试..."
	@$(GO) test -v ./...

mod-tidy:
	@echo "📦 整理依赖..."
	@$(GO) mod tidy

clean:
	@echo "🧹 清理..."
	@rm -rf bin/

help:
	@echo "可用命令:"
	@echo "  make run       - 运行服务"
	@echo "  make build     - 编译"
	@echo "  make test      - 测试"
	@echo "  make clean     - 清理"
	@echo "  make mod-tidy  - 整理依赖"
`

	return os.WriteFile(outputPath, []byte(content), 0644)
}

// generateGitignore 生成 .gitignore
func (g *DatabaseGenerator) generateGitignore() error {
	outputPath := filepath.Join(g.config.Output, ".gitignore")

	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	content := `# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test

# Output
test-output/
*.out

# Go
go.sum

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local

# Logs
*.log
`

	return os.WriteFile(outputPath, []byte(content), 0644)
}

// generateEnvExample 生成 .env.example
func (g *DatabaseGenerator) generateEnvExample() error {
	outputPath := filepath.Join(g.config.Output, ".env.example")

	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	content := `# 数据库配置
DATABASE_DSN=root:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local

# Redis 配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# 服务配置
SERVER_PORT=8080
GIN_MODE=debug

# 日志配置
LOG_LEVEL=info
`

	return os.WriteFile(outputPath, []byte(content), 0644)
}

// generateTestScript 生成测试脚本
func (g *DatabaseGenerator) generateTestScript(modelNames []string) error {
	scriptsDir := filepath.Join(g.config.Output, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(scriptsDir, "test.sh")

	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	// 生成测试脚本内容
	content := `#!/bin/bash

echo "🧪 测试 API 端点..."

BASE_URL="http://localhost:8080"

# 1. 健康检查
echo -e "\n1️⃣  健康检查"
curl -s $BASE_URL/health

`
	// 为每个模型生成测试命令
	for _, name := range modelNames {
		lowerName := toLowerCamelCaseMain(name)
		content += fmt.Sprintf(`
# %s
echo -e "\n2️⃣  获取 %s 列表"
curl -s $BASE_URL/api/v1/%s

echo -e "\n3️⃣  创建 %s"
curl -s -X POST $BASE_URL/api/v1/%s \
  -H "Content-Type: application/json" \
  -d '{"name":"test"}'

echo -e "\n4️⃣  获取 %s ID=1"
curl -s $BASE_URL/api/v1/%s/1
`, name, lowerName, lowerName, lowerName, lowerName, lowerName, lowerName)
	}

	content += `

echo -e "\n✅ 测试完成！"
`

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return err
	}

	// 设置可执行权限
	return os.Chmod(outputPath, 0755)
}

// generateDockerCompose 生成 docker-compose.yml
func (g *DatabaseGenerator) generateDockerCompose() error {
	outputPath := filepath.Join(g.config.Output, "docker-compose.yml")

	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	projectName := filepath.Base(g.config.Output)

	content := fmt.Sprintf(`version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: %s-mysql
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: mydb
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - %s-network

  redis:
    image: redis:7-alpine
    container_name: %s-redis
    ports:
      - "6379:6379"
    networks:
      - %s-network

volumes:
  mysql-data:

networks:
  %s-network:
    driver: bridge
`, projectName, projectName, projectName, projectName, projectName)

	return os.WriteFile(outputPath, []byte(content), 0644)
}
