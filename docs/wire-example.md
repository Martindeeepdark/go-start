# 依赖注入方案对比

## 当前状态（手动依赖注入）

```go
// main.go
func main() {
    db := database.New(&cfg.Database)
    cacheClient := cache.New(&cfg.Redis)

    // 手动构造依赖链
    repo := repository.New(db)
    svc := service.New(repo, cacheClient)
    ctrl := controller.New(svc)

    // 使用
    registerRoutes(r, ctrl)
}
```

**优点：**
- 简单直接
- 无需额外依赖
- 容易理解

**缺点：**
- 依赖链长时代码混乱
- 难以维护
- 测试时需要手动构建所有依赖

---

## 方案 1：使用 Wire（推荐）

### 安装 Wire
```bash
go install github.com/google/wire/cmd/wire@latest
```

### 定义 Provider

```go
// cmd/server/internal/providers/providers.go
package providers

import (
    "{{.Module}}/config"
    "{{.Module}}/internal/controller"
    "{{.Module}}/internal/repository"
    "{{.Module}}/internal/service"
    "{{.Module}}/pkg/cache"
    "{{.Module}}/pkg/database"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// ProvideConfig 提供配置
func ProvideConfig() (*config.Config, error) {
    return config.Load()
}

// ProvideLogger 提供日志
func ProvideLogger() (*zap.Logger, error) {
    return zap.NewProduction()
}

// ProvideDatabase 提供数据库
func ProvideDatabase(cfg *config.Config) (*database.DB, error) {
    return database.New(&cfg.Database)
}

// ProvideCache 提供 Redis（可选）
func ProvideCache(cfg *config.Config) (*cache.Cache, error) {
    return cache.New(&cfg.Redis)
}

// ProvideRepository 提供仓储层
func ProvideRepository(db *database.DB) *repository.Repository {
    return repository.New(db)
}

// ProvideService 提供服务层
func ProvideService(repo *repository.Repository, cache *cache.Cache) *service.Service {
    return service.New(repo, cache)
}

// ProvideController 提供控制器层
func ProvideController(svc *service.Service) *controller.Controller {
    return controller.New(svc)
}

// ProvideRouter 提供路由
func ProvideRouter(ctrl *controller.Controller, logger *zap.Logger) *gin.Engine {
    r := router.New()
    r.Use(
        httpx.Logger(logger),
        httpx.Recovery(logger),
        httpx.CORS(),
        httpx.RequestID(),
    )
    registerRoutes(r, ctrl)
    return r.Engine()
}
```

### 定义 Wire 注入器

```go
// cmd/server/internal/providers/wire.go
//go:build wireinject
// +build wireinject

package providers

import (
    "github.com/google/wire"
    "{{.Module}}/cmd/server"
)

// InitApp 初始化应用（Wire 会生成这个函数的实现）
func InitApp() (*server.App, error) {
    wire.Build(
        // 超级依赖集
        wire.Struct(new(server.App), "*"),

        // 提供者集合
        ProviderSet,
    )
    return &server.App{}, nil
}

// ProviderSet 所有提供者
var ProviderSet = wire.NewSet(
    ProvideConfig,
    ProvideLogger,
    ProvideDatabase,
    ProvideCache,
    ProvideRepository,
    ProvideService,
    ProvideController,
    ProvideRouter,
)
```

### App 结构体

```go
// cmd/server/app.go
package server

import (
    "context"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type App struct {
    Config    *config.Config
    Logger    *zap.Logger
    DB        *database.DB
    Cache     *cache.Cache
    Router    *gin.Engine
}

func (a *App) Run() error {
    addr := fmt.Sprintf(":%d", a.Config.Server.Port)
    a.Logger.Info("Starting server", zap.String("addr", addr))

    return a.Router.Run(addr)
}

func (a *App) Shutdown(ctx context.Context) error {
    a.Logger.Info("Shutting down...")
    a.DB.Close()
    if a.Cache != nil {
        a.Cache.Close()
    }
    return nil
}
```

### 生成 Wire 代码

```bash
cd cmd/server/internal/providers
wire
```

Wire 会自动生成 `wire_gen.go` 文件。

### 使用

```go
// cmd/server/main.go
func main() {
    app, err := providers.InitApp()
    if err != nil {
        log.Fatal(err)
    }
    defer app.Shutdown(context.Background())

    app.Run()
}
```

---

## 方案 2：使用 Fx（运行时依赖注入）

Fx 是 Uber 开发的依赖注入框架，运行时反射。

### 安装 Fx
```bash
go get go.uber.org/fx@latest
```

### 使用示例

```go
// cmd/server/main.go
package main

import (
    "go.uber.org/fx"
    "{{.Module}}/config"
    "{{.Module}}/internal/controller"
    "{{.Module}}/internal/repository"
    "{{.Module}}/internal/service"
)

func main() {
    fx.New(
        // 提供依赖
        fx.Provide(
            config.Load,           // 配置
            database.New,          // 数据库
            cache.New,             // 缓存
            repository.New,        // 仓储层
            service.New,           // 服务层
            controller.New,        // 控制层
            NewRouter,             // 路由
        ),

        // 启动应用
        fx.Invoke(func(r *gin.Engine) {
            r.Run(":8080")
        }),
    )

    // 生命周期管理
    fx.New(
        fx.Provide(...),
        fx.Invoke(RegisterRoutes),
        fx.Start(func() {
            // 应用启动
        }),
        fx.Stop(func() {
            // 应用关闭
        }),
    )
}
```

---

## 对比

| 特性 | Wire | Fx | 手动注入 |
|------|------|-----|---------|
| 性能 | ⭐⭐⭐⭐⭐ 编译时 | ⭐⭐⭐ 运行时反射 | ⭐⭐⭐⭐⭐ |
| 易用性 | ⭐⭐⭐⭐ 需要生成代码 | ⭐⭐⭐⭐⭐ 开箱即用 | ⭐⭐⭐ |
| 维护性 | ⭐⭐⭐⭐⭐ 清晰 | ⭐⭐⭐⭐ 较好 | ⭐⭐ 难维护 |
| 调试 | ⭐⭐⭐⭐ 编译时发现错误 | ⭐⭐⭐ 运行时发现 | ⭐⭐⭐⭐ |
| 学习曲线 | ⭐⭐⭐ 中等 | ⭐⭐⭐⭐ 简单 | ⭐⭐⭐⭐⭐ 最简单 |

---

## 推荐方案

### 对于 go-start 项目

**推荐使用 Wire**，原因：
1. ✅ 编译时生成代码，零运行时开销
2. ✅ 类型安全，编译时就能发现错误
3. ✅ 生成代码可读，方便调试
4. ✅ Google 出品，社区活跃
5. ✅ 适合大型项目

### 实施步骤

1. **第一阶段：打样**
   - 先在 User 模块实现完整的 Wire 依赖注入
   - 作为其他模块的参考

2. **第二阶段：推广**
   - 将 Wire 应用到整个项目
   - 更新文档和示例

3. **第三阶段：优化**
   - 根据使用反馈调整
   - 添加更多 Providers

---

## 示例：完整的 Wire 实现

见下一个文件：`wire-user-example.go`
