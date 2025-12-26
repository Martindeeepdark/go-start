# go-start - 高级 Go 脚手架工具

> 帮助新人工程师快速上手，让高级工程师大展身手

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Project Status](https://img.shields.io/badge/status-70%25-yellow.svg)](STATUS.md)

## ✨ 特性

- 🚀 **快速创建项目** - 一条命令生成完整的 Web 项目
- 🔄 **CRUD 自动化** - 从数据库表自动生成完整代码
- 📝 **详细中文注释** - 每个函数都有清晰的说明
- 🏗️ **多种架构** - 支持 MVC 和 DDD 架构
- 🎨 **生产级代码** - 使用 GORM Gen，类型安全

## 📦 快速开始

### 安装

```bash
go install github.com/yourname/go-start@latest
```

或从源码安装：

```bash
git clone https://github.com/yourname/go-start.git
cd go-start
go build -o go-start cmd/go-start/*.go
sudo mv go-start /usr/local/bin/
```

### 创建项目

```bash
# 创建新的 API 项目
go-start create my-api

# 指定模块名
go-start create my-api --module=github.com/username/my-api

# 使用交互式向导（TODO）
go-start create my-api --wizard
```

启动项目：

```bash
cd my-api
go mod tidy
cp config.yaml.example config.yaml
# 编辑 config.yaml 配置数据库
go run cmd/server/main.go
```

访问 http://localhost:8080/health 查看健康检查

### 生成 CRUD 代码

```bash
# 准备数据库
mysql -u root -p -e "CREATE DATABASE mydb;"

# 生成代码（交互式选择表）
go-start gen db --dsn="root:pass@tcp(localhost:3306)/mydb" --interactive

# 指定表名
go-start gen db \
  --dsn="root:pass@tcp(localhost:3306)/mydb" \
  --tables="users,articles,comments" \
  --output="./internal"
```

生成的代码包含：
- ✅ Model 层（GORM Gen）
- ✅ Repository 层（CRUD + 高级查询）
- ✅ Service 层（业务逻辑 + 缓存）
- ✅ Controller 层（RESTful API）
- ✅ 路由自动注册

## 🎯 核心功能

### 1. create 命令 - 创建项目

```bash
go-start create <project-name> [flags]
```

**选项**:
- `--arch` - 架构类型（mvc, ddd，默认 mvc）
- `--module` - Go 模块名
- `--wizard` - 使用交互式向导（TODO）

**生成的项目结构**:

```
my-api/
├── cmd/
│   └── server/
│       └── main.go          # 入口文件
├── internal/
│   ├── controller/          # 控制器层
│   ├── service/             # 业务逻辑层
│   ├── repository/          # 数据访问层
│   └── model/               # 数据模型
├── pkg/
│   ├── database/            # 数据库封装
│   ├── cache/               # Redis 缓存
│   └── httpx/               # HTTP 工具
├── config/
│   └── config.go            # 配置管理
├── go.mod
├── go.sum
├── config.yaml.example      # 配置示例
└── README.md                # 项目文档
```

### 2. gen db 命令 - 生成 CRUD 代码

```bash
go-start gen db [flags]
```

**选项**:
- `--dsn` - 数据库连接字符串（必填）
- `--tables` - 表名，逗号分隔（如：users,articles）
- `--interactive` - 交互式选择表（推荐）
- `--arch` - 架构类型（mvc, ddd，默认 mvc）
- `--output` - 输出目录（默认 ./internal）

**示例**:

```bash
# MySQL
go-start gen db \
  --dsn="root:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local" \
  --tables="users"

# PostgreSQL
go-start gen db \
  --dsn="host=localhost port=5432 user=root password=pass dbname=mydb" \
  --tables="articles"

# 使用通配符
go-start gen db --dsn="..." --tables="user*"
```

### 3. 架构选择

#### MVC 架构（默认）
```
internal/
├── controller/    # HTTP 处理器
├── service/       # 业务逻辑
├── repository/    # 数据访问
└── model/         # 数据模型
```

#### DDD 架构（TODO）
```
internal/
├── domain/           # 领域层
│   ├── entity/       # 实体
│   ├── repository/   # 仓储接口
│   └── service/      # 领域服务
├── application/      # 应用层
│   └── service/      # 应用服务
├── infrastructure/   # 基础设施层
│   └── persistence/  # 持久化实现
└── interface/        # 接口层
    └── http/         # HTTP 控制器
```

详见 [DDD 指南](docs/DDD_GUIDE.md)

## 📊 项目状态

当前进度：**70%**

| 功能 | 状态 | 完成度 |
|-----|------|--------|
| create 命令 | ✅ 可用 | 90% |
| gen db 命令 | 🟢 基本可用 | 90% |
| DDD 架构 | 🔴 未完成 | 40% |
| Spec-Kit | 🔴 未实现 | 30% |

详细状态请查看 [STATUS.md](STATUS.md)

最近更新：
- ✅ **2025-12-26**: create 命令端到端测试通过
- ✅ **2025-12-26**: gen db 命令端到端测试通过
- ✅ **2025-12-26**: 修复了多个模板 bug

## 📚 文档

### 用户文档
- [QUICKSTART.md](QUICKSTART.md) - 5 分钟快速上手
- [STATUS.md](STATUS.md) - 项目当前状态
- [ARCHITECTURE.md](ARCHITECTURE.md) - 架构设计说明

### 详细指南
- [docs/DDD_GUIDE.md](docs/DDD_GUIDE.md) - DDD 架构详细指南
- [docs/GORM_GEN_GUIDE.md](docs/GORM_GEN_GUIDE.md) - GORM Gen 使用指南
- [docs/COMPLETE_EXAMPLE.md](docs/COMPLETE_EXAMPLE.md) - 完整示例项目

### 开发文档
- [DESIGN.md](DESIGN.md) - 系统设计文档
- [TEST_RESULTS.md](TEST_RESULTS.md) - create 命令测试报告
- [GEN_DB_TEST_REPORT.md](GEN_DB_TEST_REPORT.md) - gen db 命令测试报告

文档索引请查看 [DOCS_INDEX.md](DOCS_INDEX.md)

## 🛠️ 技术栈

- **CLI**: [Cobra](https://github.com/spf13/cobra) - 命令行框架
- **Web**: [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- **ORM**: [GORM](https://github.com/go-gorm/gorm) + [GORM Gen](https://github.com/go-gorm/gen)
- **缓存**: [go-redis](https://github.com/redis/go-redis) - Redis 客户端
- **日志**: [zap](https://github.com/uber-go/zap) - 结构化日志
- **配置**: [viper](https://github.com/spf13/viper) - 配置管理
- **文档**: [swaggo](https://github.com/swaggo/gin-swagger) - Swagger 文档

## 🤝 贡献

欢迎贡献！请随时提交 Issue 或 Pull Request。

## 📄 许可证

MIT License

## 🙏 致谢

灵感来源于 [nunu](https://github.com/go-nunu/nunu) 项目

---

**注意**: 项目正在积极开发中，API 可能会有变化。建议使用稳定版本。
