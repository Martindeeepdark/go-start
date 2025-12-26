# go-start 快速开始

5 分钟从数据库表到完整的 RESTful API！

## 📋 前置条件

在开始之前，请确保你已经安装了：

- ✅ **Go 1.21 - 1.23**（推荐 1.21）
  - ⚠️ **重要**：Go 1.24+ 与 golang.org/x/tools 存在已知兼容性问题
  - 检查版本：`go version`
- ✅ **MySQL 5.7+** 或 **PostgreSQL 12+**
- ✅ **5 分钟时间** ⏰

### 快速检查 Go 版本

```bash
go version
# 输出示例：go version go1.21.0 darwin/amd64 ✅
# 输出示例：go version go1.24.0 darwin/amd64 ❌
```

**如果版本不兼容**，请先安装正确的 Go 版本。详见 [Go 版本要求](docs/VERSION_REQUIREMENTS.md)。

## 第一步：构建工具

```bash
cd /path/to/go-start
go build -o bin/go-start cmd/go-start/*.go
```

## 第二步：准备数据库

创建一个测试表：

```sql
CREATE DATABASE testdb;
USE testdb;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    age INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
);
```

## 第三步：生成代码

```bash
./bin/go-start gen db \
  --dsn="root:password@tcp(localhost:3306)/testdb" \
  --tables=users \
  --output=./myproject \
  --module=github.com/username/myproject
```

**预期输出**：

```
🔌 正在连接数据库...
📊 DSN: root:***@tcp(localhost:3306)/testdb
📋 将生成 1 张表: users
🏗️  架构模式: MVC

✅ 代码生成完成！

📦 已生成:
  ✓ Model (数据模型)
  ✓ Repository (数据访问层 + CRUD + 高级查询)
  ✓ Service (业务逻辑层 + 缓存)
  ✓ Controller (HTTP 处理器 + RESTful API)
  ✓ Routes (路由注册)
  ✓ pkg/cache (Redis 缓存封装)
  ✓ pkg/httpx/response (统一响应格式)
```

## 第四步：查看生成的代码

```bash
tree myproject -L 3
```

你会看到：

```
myproject/
├── cmd/
│   └── server/
│       └── main.go              # ✅ 应用入口（已生成）
├── internal/
│   ├── dal/                     # GORM Gen 查询 API
│   │   ├── query/
│   │   │   ├── gen.go
│   │   │   └── users.go
│   │   └── model/
│   │       └── users.gen.go
│   ├── model/
│   │   └── common.go            # ✅ 通用模型（已生成）
│   ├── repository/              # 数据访问层
│   │   └── users.go
│   ├── service/                 # 业务逻辑层（带缓存）
│   │   └── users.go
│   ├── controller/              # RESTful API
│   │   └── users.go
│   └── routes/                  # 路由注册
│       └── auto_routes.go
├── pkg/
│   ├── cache/
│   │   └── cache.go             # ✅ Redis 缓存封装（已生成）
│   └── httpx/
│       └── response/
│           └── response.go      # ✅ 统一响应格式（已生成）
├── go.mod                       # ✅ Go 模块文件（已生成）
└── config.yaml.example          # ✅ 配置文件示例（已生成）
```

## 第五步：运行服务

**✅ main.go 已经生成，无需手动编写！**

```bash
cd myproject

# 设置数据库环境变量
export DATABASE_DSN="root:password@tcp(localhost:3306)/testdb"

# 运行服务
go run cmd/server/main.go
```

**预期输出**：

```
2024/12/26 15:30:00 INFO Starting github.com/username/myproject...
2024/12/26 15:30:01 INFO Database connected successfully
2024/12/26 15:30:01 INFO [GIN-debug] Listening and serving HTTP on :8080
```

## 第六步：测试 API

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","age":25}'

# 获取列表
curl http://localhost:8080/api/v1/users

# 获取详情
curl http://localhost:8080/api/v1/users/1

# 更新
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"username":"alice_updated","age":26}'

# 删除
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 完成！🎉

你现在有了一个完整的 RESTful API，包括：

- ✅ 类型安全的数据库操作（GORM Gen）
- ✅ 数据访问层（Repository）
- ✅ 业务逻辑层（Service，带缓存支持）
- ✅ HTTP 处理层（Controller）
- ✅ 自动路由注册
- ✅ 详细的中文注释

**每个功能都有详细注释，专注于业务逻辑即可！**

## 下一步

- 查看完整示例：[docs/COMPLETE_EXAMPLE.md](./docs/COMPLETE_EXAMPLE.md)
- 了解技术选型：[docs/GORM_TECH_CHOICE.md](./docs/GORM_TECH_CHOICE.md)
- 查看功能清单：[docs/FEATURE_CHECKLIST.md](./docs/FEATURE_CHECKLIST.md)

---

**重点**：生成的所有代码都有详细中文注释，即使你是新人也能快速理解！
