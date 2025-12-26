# go-start 使用指南

## 🚀 快速开始

go-start 是一个强大的 Go 项目脚手架工具，帮助新人工程师快速上手，也让高级工程师能够大展身手！

### 安装

```bash
# 从源码安装
git clone https://github.com/yourname/go-start.git
cd go-start
go install ./cmd/go-start

# 或使用 go install
go install github.com/yourname/go-start/cmd/go-start@latest
```

---

## 🎯 两种使用模式

### 1️⃣ 传统命令行模式（快速创建）

适合：熟悉参数的高级工程师

```bash
# 使用默认配置创建项目
go-start create my-api

# 自定义模块名
go-start create my-api --module=github.com/myname/my-api

# 选择架构模式
go-start create my-api --arch=ddd
```

### 2️⃣ 交互式向导模式（推荐）⭐

适合：所有工程师，特别是新人

```bash
# 启动交互式向导
go-start create --wizard

# 或者直接运行（不带参数也会启动向导）
go-start create
```

**向导会问什么？**
1. 📦 项目名称
2. 📦 Go 模块名称
3. 📝 项目描述（可选）
4. 🏗️ 架构模式（MVC / DDD）
5. 🗄️ 数据库类型（MySQL / PostgreSQL / SQLite）
6. ⚡ 是否需要 Redis 缓存
7. 🔐 是否需要用户认证系统
8. 📚 是否需要 Swagger API 文档
9. 🔌 服务器端口

**向导特点：**
- ✨ 友好的中文界面
- 💡 每个选项都有详细说明
- 🎨 彩色输出，清晰的步骤提示
- ✅ 智能验证，防止输入错误
- 📋 配置摘要，确认后才创建

---

## 📦 创建项目示例

### 示例 1: 创建一个简单的 RESTful API

```bash
$ go-start create --wizard

╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   🚀 欢迎使用 go-start 交互式项目创建向导                  ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

📦 步骤 1/9: 项目名称
═════════════════════════════════════════
➜ 请输入项目名称（例如: my-api）: blog-api

📦 步骤 2/9: Go 模块名称
═════════════════════════════════════════
➜ 请输入 Go 模块名称 (默认: github.com/yourname/blog-api):

📝 步骤 3/9: 项目描述
═════════════════════════════════════════
➜ 请输入项目描述（可选）: 一个简单的博客API服务

🏗️ 步骤 4/9: 架构模式
═════════════════════════════════════════
选择你的项目架构模式：
  1️⃣  MVC (Model-View-Controller)
     - 适合中小型项目
     - 简单直观，易于上手
     - 推荐：新人首选

  2️⃣  DDD (Domain-Driven Design)
     - 适合大型复杂项目
     - 领域驱动设计，业务逻辑清晰
     - 推荐：高级工程师
➜ 请选择架构模式 (1 或 2) (默认: 1): 1

🗄️ 步骤 5/9: 数据库类型
═════════════════════════════════════════
选择你使用的数据库：
  1️⃣  MySQL
  2️⃣  PostgreSQL
  3️⃣  SQLite
➜ 请选择数据库 (1/2/3) (默认: 1): 1

⚡ 步骤 6/9: Redis 缓存
═════════════════════════════════════════
➜ 是否需要 Redis 支持？(y/n) (默认: y): y

🔐 步骤 7/9: 用户认证系统
═════════════════════════════════════════
➜ 是否需要认证系统？(y/n) (默认: y): y

📚 步骤 8/9: API 文档
═════════════════════════════════════════
➜ 是否需要 Swagger 文档？(y/n) (默认: y): y

🔌 步骤 9/9: 服务器端口
═════════════════════════════════════════
➜ 请输入服务器端口号 (默认: 8080): 8080

════════════════════════════════════════════════════════════
📋 项目配置摘要
════════════════════════════════════════════════════════════
  项目名称:        blog-api
  模块名称:        github.com/yourname/blog-api
  项目描述:        一个简单的博客API服务
  架构模式:        MVC (Model-View-Controller)
  数据库:          MySQL
  Redis 缓存:      ✓ 启用
  认证系统:        ✓ 启用
  Swagger 文档:    ✓ 启用
  服务端口:        8080
════════════════════════════════════════════════════════════

✨ 准备创建项目！
➜ 确认创建项目？(y/n) (default: y): y

✓ 项目创建成功！
```

### 示例 2: 使用传统命令行模式

```bash
# 一行命令创建项目
go-start create my-api --module=github.com/myname/my-api

# 项目创建完成后
cd my-api
go mod tidy
cp config.yaml.example config.yaml
# 编辑 config.yaml 配置数据库
go run cmd/server/main.go
```

---

## 📂 项目结构

```
my-api/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口点（带详细注释）
├── internal/
│   ├── controller/              # 控制器层（HTTP 处理器）
│   │   ├── controller.go        # 控制器基类
│   │   └── user.go              # 用户控制器（带详细注释）
│   ├── service/                 # 服务层（业务逻辑）
│   │   ├── service.go           # 服务接口定义
│   │   └── user.go              # 用户服务（📝 详细中文注释）
│   ├── repository/              # 数据访问层
│   │   ├── repository.go        # 仓储接口
│   │   └── user.go              # 用户仓储
│   ├── model/                   # 数据模型
│   │   └── user.go              # 用户模型
│   └── middleware/              # 中间件（如果启用认证）
│       └── auth.go              # JWT 认证中间件
├── config/
│   ├── config.go                # 配置加载
│   └── config.yaml.example      # 配置文件示例
├── pkg/                         # 私有工具库
│   ├── cache/                   # Redis 缓存封装
│   ├── database/                # 数据库封装
│   └── httpx/                   # HTTP 工具
│       ├── middleware/          # 中间件
│       ├── response/            # 统一响应格式
│       └── router/              # 路由封装
├── go.mod                       # Go 模块定义
├── go.sum
├── README.md                    # 项目文档
└── .gitignore
```

---

## 💡 代码特点

### 新人友好

**详细的中文注释：**
```go
// UserService 用户服务层
//
// 职责说明：
//   - 实现用户相关的业务逻辑
//   - 协调 repository 层和 cache 层
//   - 处理数据的缓存策略
//   - 实现业务的校验和规则
//
// 使用示例：
//   service := service.NewUserService(userRepo, cacheClient)
//   user, err := service.GetByID(ctx, 1)
type UserService struct {
    repo  *repository.UserRepository // 用户数据访问层
    cache *cache.Cache              // Redis 缓存客户端
}

// GetByID 根据 ID 获取用户信息
//
// 缓存策略：
//   1. 先从 Redis 缓存查询
//   2. 缓存命中则直接返回
//   3. 缓存未命中则查询数据库
//   4. 查询成功后写入缓存，过期时间 10 分钟
//
// 参数：
//   - ctx: 上下文
//   - id: 用户 ID
//
// 返回：
//   - *model.User: 用户对象
//   - error: 用户不存在或查询失败时返回错误
//
// 使用场景：
//   - 用户详情页
//   - 用户资料编辑
//   - 验证用户身份
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
    // 实现代码...
}
```

**每个函数都有：**
- 📝 功能说明
- 🎯 使用场景
- 📋 参数说明
- ✅ 返回值说明
- 💡 业务逻辑详解
- ⚠️ 注意事项

### 高级可扩展

**清晰的架构分层：**
```
Controller (控制器层)
    ↓
Service (服务层 - 业务逻辑)
    ↓
Repository (仓储层 - 数据访问)
    ↓
Database (数据库)
```

**易于定制：**
- 所有代码都是生成后可编辑
- 清晰的接口定义
- 灵活的依赖注入
- 可选功能模块（认证、Swagger、Redis）

---

## 🚀 启动项目

### 1. 进入项目目录

```bash
cd my-api
```

### 2. 下载依赖

```bash
go mod tidy
```

### 3. 配置数据库

```bash
cp config.yaml.example config.yaml
```

编辑 `config.yaml`：

```yaml
server:
  port: 8080

database:
  driver: mysql  # mysql, postgresql, sqlite
  host: localhost
  port: 3306
  database: my_api
  username: root
  password: ""

redis:
  host: localhost
  port: 6379
  db: 0
  password: ""
```

### 4. 运行项目

```bash
# 直接运行
go run cmd/server/main.go

# 或使用 go-start run（支持热重载）
go-start run
```

### 5. 访问服务

```bash
# 健康检查
curl http://localhost:8080/health

# 如果启用了 Swagger
# 访问 http://localhost:8080/swagger/index.html
```

---

## 🔧 开发工具

### 热重载

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 使用 go-start run（自动检测 air）
go-start run
```

### 生成 Swagger 文档

```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init

# 访问 http://localhost:8080/swagger/index.html
```

---

## 📚 更多命令

```bash
# 查看帮助
go-start help

# 查看版本
go-start version

# 运行项目（热重载）
go-start run
```

---

## 🎓 新人学习路径

### 第 1 天：理解项目结构
1. 阅读 `README.md`
2. 查看 `internal/model/user.go` - 数据模型
3. 查看 `internal/repository/user.go` - 数据访问
4. 查看 `internal/service/user.go` - 业务逻辑
5. 查看 `internal/controller/user.go` - HTTP 处理器

### 第 2 天：运行项目
1. 配置数据库
2. 运行项目
3. 测试 API
4. 查看日志
5. 理解请求流程

### 第 3 天：添加新功能
1. 创建新的 Model
2. 创建 Repository
3. 创建 Service
4. 创建 Controller
5. 注册路由

### 第 4 天：高级特性
1. 使用 Redis 缓存
2. 添加 JWT 认证
3. 生成 Swagger 文档
4. 编写单元测试

---

## 🌟 核心特性

### ✅ 新人友好
- 🎨 交互式向导，一步步引导
- 📝 详细的中文注释
- 💡 每个函数都有使用场景说明
- 🚗 清晰的代码结构和命名
- ⚠️ 友好的错误提示

### ✅ 高级可扩展
- 🏗️ 支持多种架构模式（MVC/DDD）
- 🔌 插件系统（规划中）
- 📐 规范驱动开发（规划中）
- 🎯 完全可定制的模板
- ⚡ 高性能优化

### ✅ 生产就绪
- ✨ 完整的 MVC 架构
- 🔐 JWT 认证系统
- 📚 Swagger API 文档
- ⚡ Redis 缓存支持
- 🗄️ 多数据库支持（MySQL/PostgreSQL/SQLite）
- 📊 结构化日志（zap）
- 🔄 优雅关闭
- 🎯 统一错误处理

---

## 🆚 与其他脚手架对比

| 特性 | go-start | nunu | gin-v2-admin |
|------|----------|------|--------------|
| **交互式向导** | ✅ 友好的中文向导 | ❌ | ❌ |
| **详细注释** | ✅ 每个函数都有详细说明 | ⚠️ 基础注释 | ⚠️ 基础注释 |
| **架构模式** | ✅ MVC/DDD（规划中） | ✅ MVC/DDD | ✅ MVC/DDD |
| **认证系统** | ✅ JWT（可选） | ✅ JWT | ✅ Casbin |
| **Swagger** | ✅（可选） | ⚠️ | ✅ |
| **新人友好** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| **文档质量** | 中文详细 | 英文 | 英文 |
| **spec-kit** | 🚧 规划中 | ❌ | ❌ |
| **AI 辅助** | 🚧 规划中 | ❌ | ❌ |

---

## 🔮 未来规划

### Phase 1: 基础增强 ✅
- [x] 交互式向导
- [x] 详细中文注释
- [x] 友好的错误提示

### Phase 2: 高级功能（进行中）
- [ ] spec-kit 集成
- [ ] DDD 架构模板
- [ ] 插件系统

### Phase 3: AI 增强
- [ ] Claude Skills 集成
- [ ] 智能代码生成
- [ ] 自动化测试生成

---

## 📞 获取帮助

```bash
# 查看帮助文档
go-start help

# 查看具体命令帮助
go-start create --help
go-start run --help

# 报告问题
https://github.com/yourname/go-start/issues
```

---

## 📄 许可证

MIT License

---

## 🙏 致谢

感谢使用 go-start！如果你觉得这个工具有帮助，请给我们一个 ⭐️
