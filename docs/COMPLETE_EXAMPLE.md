# 完整使用示例

本指南展示如何使用 go-start 从数据库表生成完整的 CRUD 代码。

## 准备工作

### 1. 安装 go-start

```bash
cd /path/to/go-start
go build -o bin/go-start cmd/go-start/*.go
```

### 2. 准备测试数据库

```bash
# 创建测试数据库
mysql -u root -p < test-schema.sql
```

这会创建一个名为 `testdb` 的数据库，包含以下表：
- `users` - 用户表（带 username、email、age 索引）
- `articles` - 文章表（带 author_id、published、created_at 索引）
- `tags` - 标签表（带 name 索引）
- `article_tags` - 文章标签关联表
- `comments` - 评论表（带 article_id、user_id 索引）

## 生成代码

### 方式一：交互式选择（推荐）

```bash
./bin/go-start gen db \
  --dsn="root:password@tcp(localhost:3306)/testdb" \
  --interactive
```

系统会显示所有表及其字段数、索引数，然后让你交互式选择：

```
📋 发现以下表（共 5 张）：

   [ 1] users (用户表)        8 字段   3 索引
   [ 2] articles (文章表)    8 字段   3 索引
   [ 3] tags (标签表)        3 字段   1 索引
   [ 4] article_tags (文章标签关联表)  3 字段   2 索引
   [ 5] comments (评论表)    5 字段   2 索引

📝 请选择要生成的表：
   方式：
   - 输入序号（逗号分隔）: 1,2,3
   - 输入范围: 1-5
   - 输入通配符: user*
   - 输入 all 生成所有表

👉 您的选择: 1,2
```

### 方式二：直接指定表名

```bash
./bin/go-start gen db \
  --dsn="root:password@tcp(localhost:3306)/testdb" \
  --tables=users,articles,tags \
  --output=./test-output
```

### 方式三：使用通配符

```bash
# 生成所有以 user 开头的表
./bin/go-start gen db \
  --dsn="root:password@tcp(localhost:3306)/testdb" \
  --tables="user*"

# 生成所有包含 tag 的表
./bin/go-start gen db \
  --dsn="..." \
  --tables="*tag*"
```

## 生成的代码结构

```
test-output/
└── internal/
    ├── dal/
    │   ├── query/
    │   │   ├── gen.go           # GORM Gen 主入口
    │   │   ├── users.go         # User 查询 API
    │   │   ├── articles.go      # Article 查询 API
    │   │   └── tags.go          # Tag 查询 API
    │   └── model.go             # 数据模型
    │
    ├── repository/
    │   ├── user.go              # UserRepository
    │   ├── article.go           # ArticleRepository
    │   └── tag.go               # TagRepository
    │
    ├── service/
    │   ├── user.go              # UserService (带缓存)
    │   ├── article.go           # ArticleService
    │   └── tag.go               # TagService
    │
    ├── controller/
    │   ├── user.go              # UserController (RESTful API)
    │   ├── article.go           # ArticleController
    │   └── tag.go               # TagController
    │
    └── routes/
        └── auto_routes.go       # 自动路由注册
```

## 代码特性

### 1. Repository 层

每个 Repository 都包含：

**基础 CRUD 方法：**
- `Create(ctx, model)` - 创建
- `GetByID(ctx, id)` - 根据 ID 获取
- `Update(ctx, model)` - 更新
- `Delete(ctx, id)` - 删除
- `List(ctx, page, pageSize)` - 分页查询
- `Count(ctx)` - 统计总数

**基于索引的自动生成方法（以 users 表为例）：**

```go
// 因为有 username 索引，自动生成：
ByUsername(ctx, username) (*model.User, error)
ByUsernameList(ctx, username) ([]*model.User, error)

// 因为有 email 索引，自动生成：
ByEmail(ctx, email) (*model.User, error)
ByEmailList(ctx, email) ([]*model.User, error)

// 因为有 age 索引，自动生成：
ByAge(ctx, age) (*model.User, error)
ByAgeList(ctx, age) ([]*model.User, error)
```

### 2. Service 层

每个 Service 都包含：

**业务方法：**
- `Create(ctx, model)` - 包含参数校验、唯一性检查
- `GetByID(ctx, id)` - 带缓存策略（默认 10 分钟过期）
- `Update(ctx, model)` - 包含存在性检查
- `Delete(ctx, id)` - 包含存在性检查
- `List(ctx, page, pageSize)` - 分页查询，参数校验
- `Count(ctx)` - 统计总数

**缓存策略（如果启用）：**
```go
// 读取时
1. 先从 Redis 查询
2. 缓存命中直接返回
3. 未命中查询数据库
4. 写入缓存（10分钟过期）

// 写入时
1. 执行数据库操作
2. 删除相关缓存
```

### 3. Controller 层

每个 Controller 都包含：

**RESTful API 端点：**
```go
POST   /api/v1/users       - Create          // 创建
GET    /api/v1/users       - List            // 列表
GET    /api/v1/users/:id   - GetByID         // 详情
PUT    /api/v1/users/:id   - Update          // 更新
DELETE /api/v1/users/:id   - Delete          // 删除
```

**特性：**
- 统一的响应格式（使用 `response.Success/Error`）
- 参数校验和错误处理
- Swagger 注释（可用 swaggo 生成文档）
- HTTP 状态码规范

### 4. 路由注册

自动生成 `RegisterAutoRoutes` 函数：

```go
controllers := &routes.Controllers{
    User:     controller.NewUserController(userService),
    Article:  controller.NewArticleController(articleService),
    Tag:      controller.NewTagController(tagService),
}

routes.RegisterAutoRoutes(r, controllers)
```

## 使用生成的代码

### 1. 初始化数据库连接

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "yourproject/internal/dal/query"
    "yourproject/internal/repository"
    "yourproject/internal/service"
    "yourproject/internal/controller"
    "yourproject/internal/routes"
)

func main() {
    // 连接数据库
    dsn := "root:password@tcp(localhost:3306)/testdb"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // 初始化 GORM Gen
    q := query.Use(db)

    // 初始化 Repository
    userRepo := repository.NewUserRepository(db)

    // 初始化 Service
    userService := service.NewUserService(userRepo, db, cacheClient)

    // 初始化 Controller
    userController := controller.NewUserController(userService)

    // 初始化路由
    r := gin.Default()
    controllers := &routes.Controllers{
        User: userController,
    }
    routes.RegisterAutoRoutes(r, controllers)

    // 启动服务
    r.Run(":8080")
}
```

### 2. 测试 API

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "age": 25
  }'

# 获取用户列表
curl http://localhost:8080/api/v1/users?page=1&page_size=10

# 获取用户详情
curl http://localhost:8080/api/v1/users/1

# 更新用户
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice_updated",
    "email": "alice_new@example.com",
    "age": 26
  }'

# 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 高级特性

### 1. 使用索引字段查询

```go
// 使用 Repository 层自动生成的方法
user, err := userRepo.GetByEmail(ctx, "alice@example.com")

// 或使用 Service 层（带缓存）
user, err := userService.GetByEmail(ctx, "alice@example.com")
```

### 2. 自定义业务逻辑

在生成的 Service 代码中，标记了 `TODO` 的地方可以添加自定义业务逻辑：

```go
func (s *UserService) Create(ctx context.Context, user *model.User) error {
    // TODO: 添加业务校验
    // 例如：
    // - 检查邮箱格式
    // - 检查密码强度
    // - 检查年龄范围

    // 你的自定义逻辑
    if user.Age < 18 {
        return fmt.Errorf("用户年龄必须大于18岁")
    }

    // ...
}
```

### 3. 缓存控制

默认启用缓存，可以在初始化 Service 时关闭：

```go
userService := service.NewUserService(userRepo, db, nil) // 不传 cache 即禁用
```

## 下一步

1. **检查生成的代码** - 查看 `internal/` 目录
2. **添加自定义逻辑** - 在 Service 层添加业务校验
3. **编写单元测试** - 测试 Repository 和 Service
4. **配置 Redis** - 如果需要缓存功能
5. **生成 Swagger 文档** - 使用 swaggo/swag

## 常见问题

### Q: 如何修改生成代码的模块路径？

A: 目前代码中模块路径是硬编码的 `github.com/yourname/project`，后续会支持从配置文件读取。

### Q: 如何只生成部分层的代码？

A: 可以注释掉 `pkg/gen/types.go` 中 `Generate()` 方法里不需要的层生成代码。

### Q: 如何自定义生成模板？

A: 修改 `pkg/gen/repository.go`、`service.go`、`controller.go` 中的模板常量。

### Q: 重新生成会覆盖已有代码吗？

A: 会覆盖。建议：
- 将自定义业务逻辑放在 Service 层的 `TODO` 区域
- 或者使用 `spec generate` 方式，可以更好地控制生成

---

现在你可以专注于业务逻辑，而不用手写重复的 CRUD 代码了！
