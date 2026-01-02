# 简单的依赖注入方案（不使用 Wire）

## 目标
让依赖注入更清晰、更容易理解，不引入额外的库。

## 核心原则
1. 使用接口而不是具体类型
2. 构造函数明确声明依赖
3. 在 main.go 中组装依赖链
4. 保持简单直观

---

## 改进方案

### 1. Repository 层（已经是 interface）

```go
// internal/repository/user.go
package repository

// UserRepository 接口定义
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    FindByID(ctx context.Context, id uint) (*model.User, error)
    // ...
}

// userRepository 具体实现（小写，不导出）
type userRepository struct {
    db *gorm.DB
}

// NewUserRepository 返回接口类型
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}
```

### 2. Service 层（依赖 Repository 接口）

```go
// internal/service/user.go
package service

// UserService 依赖接口
type UserService struct {
    repo  repository.UserRepository  // 接口，不是具体类型
    cache *cache.Cache
}

// NewUserService 构造函数明确声明依赖
func NewUserService(repo repository.UserRepository, cache *cache.Cache) *UserService {
    return &UserService{
        repo:  repo,
        cache: cache,
    }
}
```

### 3. Controller 层（依赖 Service）

```go
// internal/controller/user.go
package controller

// UserController 依赖 Service
type UserController struct {
    service *service.UserService
}

// NewUserController 构造函数
func NewUserService(service *service.UserService) *UserController {
    return &UserController{
        service: service,
    }
}
```

### 4. Main.go 中组装依赖（清晰的依赖链）

```go
// cmd/server/main.go
func main() {
    // 1. 加载配置
    cfg := config.Load()

    // 2. 初始化基础设施
    logger := zap.NewProduction()
    db := database.New(&cfg.Database)
    defer db.Close()

    cache := cache.New(&cfg.Redis)
    defer cache.Close()

    // 3. 依赖注入链（从下到上）
    //
    // 依赖关系：
    //   DB → Repository → Service → Controller
    //
    // 构造顺序：
    //   先构造底层，再构造上层

    // Step 1: Repository 层（依赖 DB）
    userRepo := repository.NewUserRepository(db.DB())

    // Step 2: Service 层（依赖 Repository 和 Cache）
    userService := service.NewUserService(userRepo, cache)

    // Step 3: Controller 层（依赖 Service）
    userController := controller.NewUserController(userService)

    // Step 4: 路由（依赖 Controller）
    router := setupRouter(userController)

    // 5. 启动服务器
    router.Run(":8080")
}

func setupRouter(userCtrl *controller.UserController) *gin.Engine {
    r := gin.Default()

    api := r.Group("/api/v1")
    {
        users := api.Group("/users")
        {
            users.POST("", userCtrl.Create)
            users.GET(":id", userCtrl.GetByID)
            users.PUT(":id", userCtrl.Update)
            users.DELETE(":id", userCtrl.Delete)
            users.GET("", userCtrl.List)
        }
    }

    return r
}
```

---

## 优势

### ✅ 简单直观
- 不需要学习 Wire 或 Fx
- 依赖链一目了然
- 新人容易理解

### ✅ 类型安全
- 编译时检查依赖
- IDE 自动补全支持

### ✅ 易于测试
- 可以轻松 mock 依赖

```go
// 测试示例
func TestUserService(t *testing.T) {
    // Mock Repository
    mockRepo := &MockUserRepository{
        users: make(map[uint]*model.User),
    }

    // 注入 mock 依赖
    service := service.NewUserService(mockRepo, nil)

    // 测试业务逻辑
    user := &model.User{Username: "test"}
    err := service.Create(context.Background(), user)
    assert.NoError(t, err)
}
```

### ✅ 灵活切换
- 可以替换不同的实现

```go
// 切换到不同的 Repository 实现
var userRepo repository.UserRepository

if useRedis {
    userRepo = repository.NewRedisUserRepository(redisClient)
} else {
    userRepo = repository.NewMySQLUserRepository(db)
}

service := service.NewUserService(userRepo, cache)
```

---

## 完整示例

见模板文件：
- `cmd/go-start/templates/mvc/main.go.tpl`
- `cmd/go-start/templates/mvc/repository/user.go.tpl`
- `cmd/go-start/templates/mvc/service/user.go.tpl`
- `cmd/go-start/templates/mvc/controller/user.go.tpl`

---

## 常见问题

### Q: 依赖链很长怎么办？
A: 分组初始化，使用聚合对象

```go
// 方案 1: 使用聚合对象
type Repositories struct {
    User repository.UserRepository
    // ... 其他 repositories
}

type Services struct {
    User *service.UserService
    // ... 其他 services
}

repos := &Repositories{
    User: repository.NewUserRepository(db),
}

services := &Services{
    User: service.NewUserService(repos.User, cache),
}

// 方案 2: 使用工厂函数
func NewServices(db *gorm.DB, cache *cache.Cache) *Services {
    userRepo := repository.NewUserRepository(db)

    return &Services{
        User: service.NewUserService(userRepo, cache),
    }
}
```

### Q: 如何管理循环依赖？
A: 重新设计，避免循环依赖

```go
// ❌ 错误：循环依赖
// Service A → Service B → Service A

// ✅ 正确：提取共同依赖
// Service A → Common ← Service B
```

### Q: 如何处理可选依赖？
A: 使用接口或指针

```go
// 方案 1: 接口可以为 nil
type UserService struct {
    repo  repository.UserRepository
    cache *cache.Cache  // 可以为 nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
    if s.cache != nil {
        // 使用缓存
    }
    // ...
}

// 方案 2: 使用接口
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
}

type UserService struct {
    repo  repository.UserRepository
    cache Cache  // 可以为 nil 的接口
}
```

---

## 总结

这个方案：
1. ✅ 不需要额外的库
2. ✅ 依赖关系清晰
3. ✅ 易于测试
4. ✅ 适合新手
5. ✅ 符合 Go 的习惯

**推荐：对于 go-start 项目，使用这种简单的手动依赖注入就足够了！**
