# go-start spec-kit 使用指南

## 🎯 什么是 spec-kit？

**spec-kit** 是 go-start 内置的**规范驱动代码生成工具**。它允许你用 YAML 文件定义 API 规范，然后自动生成类型安全的 Go 代码。

### 为什么使用 spec-kit？

✅ **类型安全** - 从规范生成强类型代码，减少运行时错误
✅ **可重复生成** - 规范不变，代码可以重新生成
✅ **统一规范** - 团队成员遵循相同的代码结构
✅ **减少样板代码** - CRUD、DTO、Validator 等自动生成
✅ **完全本地** - 无需远程服务，所有操作都在本地完成

---

## 🚀 快速开始

### 1. 创建规范文件

```bash
# 生成示例规范文件
go-start spec init

# 这会创建 example.spec.yaml 文件
```

### 2. 编辑规范文件

打开 `example.spec.yaml`，定义你的 API：

```yaml
spec: "1.0"
kind: API
name: BlogAPI
version: v1

project:
  module: github.com/yourname/blog-api
  author: Your Name
  description: 博客管理系统 API

# 定义数据模型
models:
  - name: Article
    table: articles
    comment: 文章表
    fields:
      - name: id
        type: uint
        primary: true
        autoIncrement: true
        comment: 文章ID

      - name: title
        type: string
        size: 200
        notNull: true
        comment: 文章标题

      - name: content
        type: text
        notNull: true
        comment: 文章内容

      # ... 更多字段
```

### 3. 验证规范文件

```bash
go-start spec validate --file=blog.spec.yaml
```

输出：
```
✅ 规范文件验证通过！

📊 规范信息:
  名称: BlogAPI
  版本: v1
  模型数量: 5
  端点数量: 14
```

### 4. 生成代码

```bash
go-start spec generate --file=blog.spec.yaml --output=./my-api
```

输出：
```
🚀 开始生成代码...

📦 生成数据模型...
  ✓ Article
  ✓ User
  ✓ Category

📦 生成数据访问层...
  ✓ ArticleRepository
  ✓ UserRepository
  ✓ CategoryRepository

📦 生成业务逻辑层...
  ✓ ArticleService
  ✓ UserService
  ✓ CategoryService

📦 生成控制器层...
  ✓ ArticleController
  ✓ UserController
  ✓ CategoryController

✅ 代码生成完成！
```

---

## 📖 规范文件详解

### 基本结构

```yaml
spec: "1.0"              # 规范版本
kind: API                # 类型（API、Model等）
name: YourAPI            # API 名称
version: v1              # 版本号

project:                 # 项目配置
  module: github.com/yourname/project
  author: Your Name
  description: Project description
```

### 定义数据模型

```yaml
models:
  - name: User           # 模型名称（PascalCase）
    table: users         # 数据库表名
    comment: 用户表      # 注释说明
    fields:             # 字段定义
      - name: id
        type: uint
        primary: true
        autoIncrement: true
        comment: 用户ID

      - name: username
        type: string
        size: 50
        notNull: true
        unique: true
        comment: 用户名

      - name: email
        type: string
        size: 100
        notNull: true
        unique: true
        comment: 邮箱

      - name: status
        type: int
        default: 1
        comment: 状态 1-正常 0-禁用

      - name: created_at
        type: timestamp
        autoCreateTime: true
        comment: 创建时间

      - name: updated_at
        type: timestamp
        autoUpdateTime: true
        comment: 更新时间
```

#### 字段类型支持

| 类型 | Go 类型 | 说明 |
|------|---------|------|
| `uint` | `uint` | 无符号整数 |
| `int` | `int` | 整数 |
| `string` | `string` | 字符串 |
| `text` | `string` | 长文本 |
| `bool` | `bool` | 布尔值 |
| `float` | `float64` | 浮点数 |
| `timestamp` | `time.Time` | 时间戳 |
| `date` | `time.Time` | 日期 |
| `datetime` | `time.Time` | 日期时间 |

#### 字段属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `primary` | bool | 是否为主键 |
| `autoIncrement` | bool | 是否自增 |
| `notNull` | bool | 是否非空 |
| `unique` | bool | 是否唯一 |
| `index` | bool | 是否索引 |
| `size` | int | 字段大小 |
| `default` | string | 默认值 |
| `foreignKey` | string | 外键（如：users.id） |
| `autoCreateTime` | bool | 自动创建时间 |
| `autoUpdateTime` | bool | 自动更新时间 |
| `json` | string | 自定义 JSON tag |
| `comment` | string | 注释 |

### 定义 API 端点

```yaml
endpoints:
  - method: POST
    path: /articles
    handler: CreateArticle
    auth: true
    permission: article.create
    validate: CreateArticleRequest
    comment: 创建文章

  - method: GET
    path: /articles
    handler: ListArticles
    auth: false
    cache:
      enabled: true
      ttl: 300
    pagination:
      page: 1
      pageSize: 20
      maxPageSize: 100
    comment: 获取文章列表

  - method: GET
    path: /articles/:id
    handler: GetArticle
    auth: false
    cache:
      enabled: true
      ttl: 600
    comment: 获取文章详情
```

#### 端点属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `method` | string | HTTP 方法（GET/POST/PUT/DELETE/PATCH） |
| `path` | string | 路径（支持 :id 参数） |
| `handler` | string | 处理器名称 |
| `auth` | bool | 是否需要认证 |
| `permission` | string | 权限标识 |
| `validate` | string | 请求验证器名称 |
| `cache` | object | 缓存配置 |
| `pagination` | object | 分页配置 |
| `comment` | string | 注释说明 |

### 定义请求验证

```yaml
requests:
  - name: CreateArticleRequest
    comment: 创建文章请求
    fields:
      - name: title
        rules: required,min=5,max=200
        comment: 文章标题

      - name: content
        rules: required,min=10
        comment: 文章内容

      - name: category_id
        rules: required,numeric
        comment: 分类ID
```

---

## 📂 生成的代码结构

```
my-api/
├── internal/
│   ├── model/              # 数据模型
│   │   ├── user.go
│   │   ├── article.go
│   │   └── category.go
│   ├── repository/         # 数据访问层
│   │   ├── user.go
│   │   ├── article.go
│   │   └── category.go
│   ├── service/            # 业务逻辑层
│   │   ├── user.go
│   │   ├── article.go
│   │   └── category.go
│   └── controller/         # 控制器层
│       ├── user.go
│       ├── article.go
│       └── category.go
└── internal/routes/
    └── auto_routes.go      # 自动生成的路由注册
```

---

## 💡 使用示例

### 示例 1: 创建一个简单的博客 API

**步骤 1: 创建规范文件**

```yaml
# blog.spec.yaml
spec: "1.0"
kind: API
name: BlogAPI
version: v1

project:
  module: github.com/yourname/blog-api
  description: 简单的博客 API

models:
  - name: Article
    table: articles
    comment: 文章表
    fields:
      - name: id
        type: uint
        primary: true
        autoIncrement: true

      - name: title
        type: string
        size: 200
        notNull: true

      - name: content
        type: text
        notNull: true

      - name: created_at
        type: timestamp
        autoCreateTime: true

endpoints:
  - method: POST
    path: /articles
    handler: CreateArticle
    auth: true

  - method: GET
    path: /articles
    handler: ListArticles
    auth: false
    pagination: true

  - method: GET
    path: /articles/:id
    handler: GetArticle
    auth: false

  - method: PUT
    path: /articles/:id
    handler: UpdateArticle
    auth: true

  - method: DELETE
    path: /articles/:id
    handler: DeleteArticle
    auth: true
```

**步骤 2: 生成代码**

```bash
go-start spec generate --file=blog.spec.yaml --output=./blog-api
```

**步骤 3: 查看生成的代码**

```bash
cd blog-api
ls -la internal/model/
# article.go  - Article 数据模型

ls -la internal/service/
# article.go  - ArticleService 业务逻辑层

ls -la internal/controller/
# article.go  - ArticleController HTTP 处理器
```

**步骤 4: 集成到项目**

```go
// main.go
package main

import (
    "github.com/yourname/blog-api/internal/controller"
    "github.com/yourname/blog-api/internal/repository"
    "github.com/yourname/blog-api/internal/service"
    "github.com/yourname/blog-api/internal/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 初始化各层
    db := initDB()  // 初始化数据库

    articleRepo := repository.NewArticleRepository(db)
    articleService := service.NewArticleService(articleRepo, nil)
    articleController := controller.NewArticleController(articleService)

    // 注册控制器
    controllers := &controller.Controllers{
        Article: articleController,
    }

    // 注册路由
    routes.RegisterAutoRoutes(r, controllers)

    r.Run(":8080")
}
```

### 示例 2: 带关联的模型

```yaml
models:
  - name: User
    table: users
    fields:
      - name: id
        type: uint
        primary: true
        autoIncrement: true

      - name: username
        type: string
        size: 50
        notNull: true
        unique: true

  - name: Article
    table: articles
    fields:
      - name: id
        type: uint
        primary: true
        autoIncrement: true

      - name: title
        type: string
        size: 200
        notNull: true

      - name: author_id
        type: uint
        notNull: true
        foreignKey: users.id  # 外键关联
```

---

## 🔧 高级功能

### 批量生成

```bash
# 从目录批量生成所有规范文件
go-start spec generate --dir=./specs --output=./my-api
```

### 缓存配置

```yaml
endpoints:
  - method: GET
    path: /articles/:id
    handler: GetArticle
    cache:
      enabled: true
      ttl: 600        # 缓存10分钟
```

生成的代码会自动处理缓存逻辑。

### 分页配置

```yaml
endpoints:
  - method: GET
    path: /articles
    handler: ListArticles
    pagination:
      page: 1
      pageSize: 20
      maxPageSize: 100
```

生成的控制器会自动处理分页参数。

---

## 🎓 最佳实践

### 1. 规范文件组织

```
project/
├── specs/
│   ├── user.spec.yaml
│   ├── article.spec.yaml
│   ├── comment.spec.yaml
│   └── category.spec.yaml
└── cmd/
    └── main.go
```

### 2. 规范命名约定

- 文件名：`{feature}.spec.yaml`
- 模型名：PascalCase（如 `User`、`Article`）
- 表名：snake_case（如 `users`、`articles`）
- 字段名：snake_case（如 `created_at`）

### 3. 分模块管理

按功能模块划分规范文件，便于维护：

```yaml
# user.spec.yaml - 用户模块
spec: "1.0"
kind: API
name: UserAPI
models:
  - name: User
    fields: ...

# article.spec.yaml - 文章模块
spec: "1.0"
kind: API
name: ArticleAPI
models:
  - name: Article
    fields: ...
```

---

## 🆚 对比手动编写代码

### 手动编写（传统方式）

```go
// 需要手动编写每个字段
type User struct {
    ID       uint   `gorm:"primarykey" json:"id"`
    Username string `gorm:"size:50;not null" json:"username"`
    Email    string `gorm:"size:100;not null" json:"email"`
    // ... 手动定义所有字段
}

// 需要手动编写 CRUD 方法
func (r *UserRepository) Create(ctx context.Context, user *User) error {
    // ... 手动实现
}
// ... 需要手动实现所有方法
```

### 使用 spec-kit（规范驱动）

```yaml
# 只需定义规范
models:
  - name: User
    fields:
      - name: id
        type: uint
        primary: true
      - name: username
        type: string
        size: 50
        notNull: true
# ... 一条命令生成所有代码
```

**优势：**
- ✅ 减少 80% 的样板代码
- ✅ 统一的代码风格
- ✅ 类型安全
- ✅ 可重复生成
- ✅ 易于重构

---

## ❓ 常见问题

### Q1: 生成的代码可以修改吗？

**A:** 可以！生成的代码是标准 Go 代码，完全可编辑。只需注意：如果重新生成，会覆盖你的修改。

**建议：**
- 生成基础代码
- 在生成的代码基础上添加业务逻辑
- 保留规范文件，需要时重新生成

### Q2: 如何处理复杂的业务逻辑？

**A:** spec-kit 生成基础 CRUD 代码，复杂业务逻辑建议：
1. 在 Service 层添加自定义方法
2. 使用继承或组合扩展功能
3. 保持规范文件简单，手动编写复杂逻辑

### Q3: 如何处理数据库迁移？

**A:** 生成的 Model 可以配合 GORM AutoMigrate 或迁移工具：

```go
db.AutoMigrate(&model.User{}, &model.Article{})
```

或使用 golang-migrate 等工具。

### Q4: 支持关联关系吗？

**A:** 支持！通过 `foreignKey` 定义：

```yaml
- name: author_id
  type: uint
  foreignKey: users.id
```

生成的代码会包含关联结构。

---

## 📚 更多资源

- **完整示例**: `spec/example.blog.spec.yaml`
- **交互式向导**: `go-start create --wizard`
- **项目文档**: `README.md`
- **向导文档**: `WIZARD.md`

---

## 🚀 下一步

1. ✅ 运行 `go-start spec init` 查看示例
2. ✅ 编辑规范文件，定义你的 API
3. ✅ 运行 `go-start spec generate` 生成代码
4. ✅ 查看生成的代码，理解结构
5. ✅ 在生成的代码基础上添加业务逻辑
6. ✅ 运行 `go mod tidy` 并测试

**祝你使用愉快！** 🎉
