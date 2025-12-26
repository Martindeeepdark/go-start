package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateSupportPackages 生成支持包（cache, model, response）
func (g *DatabaseGenerator) GenerateSupportPackages() error {
	modulePath := getModulePath(g.config.Module)

	fmt.Println("📦 正在生成支持包...")

	// 1. pkg/cache
	if err := g.generateCachePackage(modulePath); err != nil {
		return err
	}

	// 2. internal/dal/model
	if err := g.generateDalModelPackage(modulePath); err != nil {
		return err
	}

	// 3. internal/model
	if err := g.generateModelPackage(modulePath); err != nil {
		return err
	}

	// 4. pkg/httpx/response
	if err := g.generateResponsePackage(modulePath); err != nil {
		return err
	}

	fmt.Println("     ✓ 所有支持包创建成功")
	return nil
}

// generateCachePackage 生成 cache 包
func (g *DatabaseGenerator) generateCachePackage(modulePath string) error {
	outputDir := filepath.Join(g.config.Output, "pkg", "cache")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, "cache.go")

	content := fmt.Sprintf(`package cache

import (
	"context"
	"time"
)

// Cache 简单的缓存接口
type Cache struct {
	// TODO: 实现 Redis 缓存
	// 这里提供一个简单实现，实际使用时可以替换为 Redis
}

// New 创建一个新的缓存实例
func New() *Cache {
	return &Cache{}
}

// Set 设置缓存
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// TODO: 实现 Redis SET
	return nil
}

// Get 获取缓存
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	// TODO: 实现 Redis GET
	return "", nil
}

// Del 删除缓存
func (c *Cache) Del(ctx context.Context, keys ...string) error {
	// TODO: 实现 Redis DEL
	return nil
}

// Close 关闭缓存连接
func (c *Cache) Close() error {
	// TODO: 关闭 Redis 连接
	return nil
}
`)

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 cache.go 失败: %w", err)
	}

	return nil
}

// generateDalModelPackage 生成 dal/model 包
func (g *DatabaseGenerator) generateDalModelPackage(modulePath string) error {
	outputDir := filepath.Join(g.config.Output, "internal", "dal", "model")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, "model.go")

	content := fmt.Sprintf(`package model

// 通用数据模型
// 这个包导出生成的模型，供其他层使用
`)

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 dal/model.go 失败: %w", err)
	}

	return nil
}

// generateModelPackage 生成 internal/model 包
func (g *DatabaseGenerator) generateModelPackage(modulePath string) error {
	outputDir := filepath.Join(g.config.Output, "internal", "model")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, "common.go")

	content := `package model

import (
	"time"
)

// Response 通用响应结构
type Response struct {
	Code    int         ` + "`json:\"code\"`" + `
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
}

// PageRequest 分页请求
type PageRequest struct {
	Page     int ` + "`json:\"page\" form:\"page\"`" + `
	PageSize int ` + "`json:\"page_size\" form:\"page_size\"`" + `
}

// PageResponse 分页响应
type PageResponse struct {
	Total int64       ` + "`json:\"total\"`" + `
	List  interface{} ` + "`json:\"list\"`" + `
	Page  int         ` + "`json:\"page\"`" + `
	Size  int         ` + "`json:\"size\"`" + `
}

// BaseModel 基础模型
type BaseModel struct {
	ID        uint       ` + "`json:\"id\"`" + `
	CreatedAt time.Time  ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time  ` + "`json:\"updated_at\"`" + `
	DeletedAt *time.Time ` + "`json:\"deleted_at,omitempty\"`" + `
}
`

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 common.go 失败: %w", err)
	}

	return nil
}

// generateResponsePackage 生成 pkg/httpx/response 包
func (g *DatabaseGenerator) generateResponsePackage(modulePath string) error {
	outputDir := filepath.Join(g.config.Output, "pkg", "httpx", "response")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, "response.go")

	content := fmt.Sprintf(`package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, gin.H{
		"code":    -1,
		"message": message,
	})
}

// ErrorWithCode 返回带错误码的响应
func ErrorWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
	})
}

// Page 返回分页响应
func Page(c *gin.Context, total int64, list interface{}, page, size int) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total": total,
			"list":  list,
			"page":  page,
			"size":  size,
		},
	})
}
`)

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 response.go 失败: %w", err)
	}

	return nil
}
