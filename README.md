# go-start

<div align="center">

**🚀 5 分钟从数据库到可用 API**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

</div>

---

## 📖 简介

**go-start** 是一个数据库驱动的 Go API 代码生成器，帮助你在 5 分钟内从现有数据库生成完整的 CRUD API。

### 核心特性

- ✅ **自动生成完整分层代码**：Model、Repository、Service、Controller、Routes
- 🔒 **类型安全**：基于 GORM Gen，编译时检查，IDE 自动补全
- ⚡ **索引查询自动生成**：基于数据库索引自动生成高效查询方法
- 💾 **内置缓存支持**：Service 层自动集成 Redis 缓存
- 🏗️ **灵活的架构**：支持 MVC 和 DDD 两种架构模式
- 🎯 **开箱即用**：自动生成 `main.go` 和 `go.mod`，无需手动配置

---

## 🎯 适用场景

| 你的需求 | go-start | 其他工具 |
|---------|----------|---------|
| 有数据库设计，想快速生成 API | ✅ 完美适配 | ❌ 需要手动编写 |
| 想要类型安全的查询 API | ✅ GORM Gen | 🟡 运行时魔法字符串 |
| 需要快速迭代和原型开发 | ✅ 重新生成即可 | ❌ 手动维护成本高 |
| 新人快速上手 Go Web 开发 | ✅ 详细中文注释 | 🟡 需要理解架构 |
| 高级工程师的 DDD 架构 | ✅ 即插即用 | 🟡 需要手动搭建 |

---

## 🚀 快速开始

### 1 分钟体验

```bash
# 安装
go install github.com/yourname/go-start@latest

# 从数据库生成代码
go-start gen db \
  --dsn="root:pass@tcp(localhost:3306)/mydb" \
  --tables=users,posts \
  --module=github.com/username/my-api

# 运行
cd my-api
export DATABASE_DSN="root:pass@tcp(localhost:3306)/mydb"
go run cmd/server/main.go

# 测试 API
curl http://localhost:8080/api/v1/users
```

### 详细教程

- 📚 **[5 分钟快速开始](QUICKSTART.md)** - 新手必读
- 📖 **[详细教程](docs/TUTORIAL.md)** - 进阶功能
- ⚙️ **[配置参考](docs/CONFIGURATION.md)** - 完整参数说明

---

## ✨ 核心功能

### 🎨 自动生成完整分层架构

```
internal/
├── dal/              # GORM Gen 生成的类型安全查询 API
│   ├── query/
│   └── model/
├── repository/       # 数据访问层（CRUD + 索引查询）
├── service/          # 业务逻辑层（内置缓存）
├── controller/       # HTTP 处理层（RESTful API）
└── routes/           # 路由自动注册
```

### 🔒 类型安全的查询 API

```go
// ✅ 类型安全，IDE 自动补全，编译时检查
user, err := r.q.User.WithContext(ctx).
    Where(r.q.User.Username.Eq("alice")).  // 无魔法字符串
    Where(r.q.User.Age.Gte(18)).
    First()
```

### ⚡ 索引查询自动生成

```go
// 如果有 idx_username 索引，自动生成
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
    return r.q.User.WithContext(ctx).
        Where(r.q.User.Username.Eq(username)).
        First()
}
```

### 💾 内置缓存支持

```go
// Service 层自动包含缓存逻辑
user, err := userService.GetByID(ctx, 1)
// 首次查询数据库，后续从缓存读取（10 分钟过期）
```

---

## 📊 与其他工具对比

### vs nunu

| 特性 | go-start | nunu |
|------|----------|------|
| **核心理念** | 数据库驱动生成 | 项目脚手架 |
| **从数据库生成** | ✅ 1 分钟生成完整 CRUD | ❌ 需要手动编写 |
| **类型安全** | ✅ GORM Gen (编译时) | 🟡 GORM (运行时) |
| **架构选择** | ✅ MVC + DDD | 🟡 固定架构 |
| **依赖注入** | ✅ 手动 + 可选 Wire | ✅ Wire |
| **学习曲线** | ✅ 新人友好 | 🟡 需要理解架构 |

**使用建议**：
- ✅ **go-start**：有数据库设计，想快速生成 API
- ✅ **nunu**：从零开始搭建项目架构

### vs 手动开发

| 对比项 | go-start | 手动开发 |
|--------|----------|---------|
| 开发时间 | 5 分钟 | 2-3 天 |
| 代码质量 | 生产级 | 因人而异 |
| 维护成本 | 低（重新生成） | 高（手动维护） |
| 类型安全 | ✅ | ❌ |
| 最佳实践 | ✅ 自动遵循 | 需要经验 |

---

## 🛠️ 安装

### 系统要求

- ✅ **Go 1.21 - 1.23** （推荐 1.21）
  - ⚠️ **注意**：Go 1.24+ 与 golang.org/x/tools 存在已知兼容性问题
  - 详见：[Go 版本要求](QUICKSTART.md)
- ✅ **MySQL 5.7+** 或 **PostgreSQL 12+**

### 检查 Go 版本

\`\`\`bash
go version
# 输出示例：go version go1.21.0 darwin/amd64 ✅
# 输出示例：go version go1.24.0 darwin/amd64 ❌
\`\`\`

**如果版本不兼容**，请先安装正确的 Go 版本。

### 从源码安装

```bash
git clone https://github.com/yourname/go-start.git
cd go-start
go build -o bin/go-start ./cmd/go-start/
sudo mv bin/go-start /usr/local/bin/
```

### 使用 go install

```bash
go install github.com/yourname/go-start@latest
```

### 验证安装

```bash
go-start --version
go-start gen db --help
```

---

## 📖 使用示例

### 基础用法

```bash
# 从 MySQL 生成代码
go-start gen db \
  --dsn="root:pass@tcp(localhost:3306)/mydb" \
  --tables=users,posts,comments

# 从 PostgreSQL 生成代码
go-start gen db \
  --dsn="host=localhost user=root password=secret dbname=mydb" \
  --tables=users

# 使用通配符
go-start gen db \
  --dsn="..." \
  --tables="user*"

# 交互式选择表（推荐）
go-start gen db \
  --dsn="..." \
  --interactive
```

### 高级用法

```bash
# 指定模块路径
go-start gen db \
  --dsn="..." \
  --tables=users \
  --module=github.com/username/my-api

# 使用 DDD 架构
go-start gen db \
  --dsn="..." \
  --tables=users \
  --arch=ddd

# 指定输出目录
go-start gen db \
  --dsn="..." \
  --tables=users \
  --output=./my-api/internal
```

---

## 🏗️ 生成的代码结构

### MVC 架构（默认）

```
my-api/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口（自动生成）
├── internal/
│   ├── dal/                     # GORM Gen 查询 API
│   │   ├── query/
│   │   │   ├── gen.go
│   │   │   └── users.go
│   │   └── model/
│   │       └── users.gen.go
│   ├── repository/              # 数据访问层
│   │   └── users.go
│   ├── service/                 # 业务逻辑层
│   │   └── users.go
│   ├── controller/              # HTTP 处理层
│   │   └── users.go
│   ├── routes/                  # 路由注册
│   │   └── auto_routes.go
│   └── model/                   # 领域模型
│       └── user.go
├── config/
│   └── config.yaml.example      # 配置文件示例
├── go.mod                       # Go 模块文件（自动生成）
└── README.md                    # 项目说明
```

### DDD 架构（可选）

```
my-api/
├── internal/
│   ├── domain/                  # 领域层
│   │   └── user/
│   │       ├── User.go          # 实体
│   │       ├── repository.go    # 仓储接口
│   │       └── service.go       # 领域服务
│   ├── application/             # 应用层
│   │   └── user/
│   │       └── service.go       # 应用服务
│   ├── infrastructure/          # 基础设施层
│   │   └── persistence/
│   │       └── UserRepositoryImpl.go
│   └── interface/               # 接口层
│       └── http/
│           └── user/
│               └── controller.go
```

---

## 🔧 配置说明

### 环境变量

生成的项目支持以下环境变量：

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `DATABASE_DSN` | 数据库连接字符串 | `root:pass@tcp(localhost:3306)/mydb` |
| `REDIS_ADDR` | Redis 地址 | `localhost:6379` |
| `SERVER_PORT` | 服务器端口 | `8080` |

### 数据库支持

- ✅ MySQL 5.7+
- ✅ PostgreSQL 12+
- 🟡 SQLite 3+（计划中）
- 🟡 MongoDB（计划中）

---

## 📚 文档

- 📖 **[5 分钟快速开始](QUICKSTART.md)** - 新手必读教程
- 🎓 **[详细教程](docs/TUTORIAL.md)** - 进阶功能和最佳实践
- ⚙️ **[配置参考](docs/CONFIGURATION.md)** - 完整参数说明
- 🏗️ **[架构设计](docs/ARCHITECTURE.md)** - 技术架构详解
- ❓ **[常见问题](docs/FAQ.md)** - 问题解答

---

## 🤝 参与贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

详见 [贡献指南](CONTRIBUTING.md)

---

## 📝 路线图

### v0.1.0 (当前版本 - MVP)
- ✅ gen db 命令（从数据库生成代码）
- ✅ MVC 架构支持
- ✅ 类型安全查询（GORM Gen）
- ✅ 索引查询自动生成
- ✅ 自动生成 main.go 和 go.mod

### v0.2.0 (计划中)
- 🔨 DDD 架构完善
- 🔨 交互式向导优化
- 🔨 单元测试模板生成
- 🔨 错误提示优化
- 🔨 进度条显示

### v0.3.0 (规划中)
- 🔮 Spec-Kit 支持（从 YAML 规范生成）
- 🔮 代码增量更新（不覆盖自定义代码）
- 🔮 Wire 依赖注入集成
- 🔮 Swagger 文档自动生成
- 🔮 SQLite 和 MongoDB 支持

---

## 🙏 致谢

- [GORM](https://github.com/go-gorm/gorm) - 强大的 Go ORM 库
- [GORM Gen](https://github.com/go-gorm/gen) - 类型安全的 DAO 生成器
- [Gin](https://github.com/gin-gonic/gin) - 高性能 Go Web 框架
- [nunu](https://github.com/go-nunu/nunu) - 优秀的 Go 项目脚手架工具
- [Cobra](https://github.com/spf13/cobra) - 强大的 CLI 应用框架

---

## 📄 开源协议

本项目基于 [MIT License](LICENSE) 开源。

---

## 📮 联系方式

- **问题反馈**: [GitHub Issues](https://github.com/yourname/go-start/issues)
- **功能建议**: [GitHub Discussions](https://github.com/yourname/go-start/discussions)
- **邮件**: yourname@example.com

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！**

Made with ❤️ by [Your Name](https://github.com/yourname)

</div>
