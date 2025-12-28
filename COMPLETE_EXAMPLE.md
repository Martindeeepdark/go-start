# 完整示例:从零到运行

这个文档会带你从零开始,一步步创建一个完整的 go-start 项目并运行起来。

## 📋 准备工作

### 1. 检查 Go 版本

```bash
go version
# 推荐输出: go version go1.21.x darwin/amd64
# ⚠️  如果是 go1.24+,可能会有兼容性问题
```

### 2. 安装 MySQL

**macOS:**
```bash
brew install mysql
brew services start mysql
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
```

**Docker (推荐):**
```bash
docker run --name mysql-test \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=testdb \
  -p 3306:3306 \
  -d mysql:8.0
```

### 3. 创建测试数据库和表

```sql
-- 登录 MySQL
mysql -u root -p123456

-- 创建数据库
CREATE DATABASE testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE testdb;

-- 创建用户表
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    age INT DEFAULT 0 COMMENT '年龄',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入测试数据
INSERT INTO users (username, email, age) VALUES
('alice', 'alice@example.com', 25),
('bob', 'bob@example.com', 30),
('charlie', 'charlie@example.com', 28);

-- 验证数据
SELECT * FROM users;
```

## 🚀 创建项目

### 方式 1: 使用向导 (推荐新手)

```bash
# 1. 运行向导
go-start create --wizard

# 2. 按提示输入:
#    - 项目名称: my-api
#    - 项目描述: 我的第一个 API
#    - 架构模式: mvc (推荐新手)
#    - 数据库: mysql
#    - 端口: 8080 (直接回车使用默认)
#    - 是否启用 Redis: n (新手可以先不启用)

# 3. 进入项目目录
cd my-api

# 4. 配置数据库连接
cp config.yaml.example config.yaml
# 编辑 config.yaml,修改数据库密码
```

### 方式 2: 使用命令行 (推荐有经验者)

```bash
# 1. 创建项目
go-start create my-api \
  --arch mvc \
  --database mysql \
  --port 8080 \
  --description "我的第一个 API"

# 2. 进入项目目录
cd my-api

# 3. 配置数据库连接
cp config.yaml.example config.yaml
# 编辑 config.yaml,修改数据库密码
```

## ⚙️ 配置数据库

编辑 `config.yaml`:

```yaml
server:
  port: 8080

database:
  driver: mysql
  host: localhost
  port: 3306
  database: testdb        # 修改为你的数据库名
  username: root          # 修改为你的数据库用户
  password: "123456"      # ⚠️ 修改为你的数据库密码
  charset: utf8mb4
  parse_time: true
  loc: Local
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
  log_level: info
```

## 🔨 生成代码

### 从数据库生成 CRUD API

```bash
# 方式 1: 使用 DSN (简单)
go-start gen db \
  --dsn="root:123456@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local" \
  --tables=users

# 方式 2: 使用配置文件 (推荐)
# 确保已配置 config.yaml
go-start gen db \
  --config=config.yaml \
  --tables=users

# 方式 3: 交互式选择表 (最直观)
go-start gen db \
  --dsn="root:123456@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local" \
  --interactive
```

### 生成的文件结构

```
my-api/
├── cmd/
│   └── server/
│       └── main.go              # ✅ 应用入口
├── internal/
│   ├── dal/                     # ✅ GORM Gen 查询 API
│   │   ├── query/
│   │   │   ├── gen.go
│   │   │   └── users.go
│   │   └── model/
│   │       └── users.gen.go
│   ├── repository/              # ✅ 数据访问层
│   │   └── users.go
│   ├── service/                 # ✅ 业务逻辑层
│   │   └── users.go
│   ├── controller/              # ✅ HTTP 处理层
│   │   └── users.go
│   └── routes/                  # ✅ 路由注册
│       └── auto_routes.go
├── config/
│   └── config.yaml.example      # ✅ 配置文件示例
├── go.mod                       # ✅ Go 模块文件
└── README.md                    # ✅ 项目说明
```

## 🏃 运行项目

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行服务器

```bash
# 开发模式 (推荐)
go run cmd/server/main.go

# 或编译后运行
go build -o bin/server cmd/server/main.go
./bin/server
```

### 3. 验证运行

服务器启动后,你应该看到类似的输出:

```
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8080
```

## 🧪 测试 API

### 1. 获取所有用户

```bash
curl http://localhost:8080/api/v1/users
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "username": "alice",
      "email": "alice@example.com",
      "age": 25,
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "username": "bob",
      "email": "bob@example.com",
      "age": 30,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 2. 获取单个用户

```bash
curl http://localhost:8080/api/v1/users/1
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "age": 25,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "david",
    "email": "david@example.com",
    "age": 35
  }'
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 4,
    "username": "david",
    "email": "david@example.com",
    "age": 35,
    "created_at": "2024-01-01T12:34:56Z"
  }
}
```

### 4. 更新用户

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice_updated",
    "email": "alice_new@example.com",
    "age": 26
  }'
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice_updated",
    "email": "alice_new@example.com",
    "age": 26,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 5. 删除用户

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success"
}
```

## 🔍 常见问题

### 问题 1: 数据库连接失败

**错误信息:**
```
Error 1045: Access denied for user 'root'@'localhost'
```

**解决方法:**
1. 检查 `config.yaml` 中的用户名和密码是否正确
2. 确认 MySQL 服务正在运行: `mysql -u root -p`
3. 检查数据库是否已创建: `SHOW DATABASES;`

### 问题 2: 端口被占用

**错误信息:**
```
bind: address already in use
```

**解决方法:**
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>

# 或修改 config.yaml 中的端口
server:
  port: 9000  # 改为其他端口
```

### 问题 3: Go 版本不兼容

**错误信息:**
```
type func(i *Charset, j *Charset) bool does not match inferred type func(a *Charset, b *Charset) int
```

**解决方法:**
```bash
# 降级到 Go 1.21-1.23
brew install go@1.21

# 或使用环境变量
GOTOOLCHAIN=local go1.21 start gen db --dsn="..."
```

### 问题 4: 找不到生成的代码

**问题:**
```
Error: cannot find package "xxx/internal/query"
```

**解决方法:**
```bash
# 1. 确保已运行 gen db 命令
go-start gen db --dsn="..." --tables=users

# 2. 检查生成的文件是否存在
ls -la internal/dal/query/

# 3. 重新生成代码
rm -rf internal/dal
go-start gen db --dsn="..." --tables=users
```

## 📚 下一步

现在你已经成功创建并运行了一个 go-start 项目!

**推荐学习路径:**

1. **理解代码结构** → 阅读 [架构设计文档](docs/ARCHITECTURE.md)
2. **添加更多功能** → 查看 [详细教程](docs/TUTORIAL.md)
3. **部署到生产** → 参考 [部署指南](docs/DEPLOYMENT.md)
4. **最佳实践** → 阅读 [最佳实践](docs/BEST_PRACTICES.md)

## 🆘 获取帮助

- 📖 **文档**: [docs/](docs/)
- 🐛 **问题反馈**: [GitHub Issues](https://github.com/Martindeeepdark/go-start/issues)
- 💬 **讨论**: [GitHub Discussions](https://github.com/Martindeeepdark/go-start/discussions)

---

**恭喜!** 🎉 你已经完成了从零到运行的完整流程!
