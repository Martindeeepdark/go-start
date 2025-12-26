# go-start 快速开始

5 分钟从数据库表到完整的 RESTful API！

## 前提条件

- Go 1.25+
- MySQL 或 PostgreSQL
- （可选）Redis

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
  --output=./myproject
```

## 第四步：查看生成的代码

```bash
tree myproject/internal
```

你会看到：

```
myproject/internal/
├── dal/
│   ├── query/
│   │   ├── gen.go      # GORM Gen API
│   │   └── users.go
│   └── model.go
├── repository/
│   └── user.go         # 数据访问层
├── service/
│   └── user.go         # 业务逻辑层（带缓存）
├── controller/
│   └── user.go         # RESTful API
└── routes/
    └── auto_routes.go  # 路由注册
```

## 第五步：初始化项目

创建 `main.go`:

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "yourmodule/internal/dal/query"
    "yourmodule/internal/repository"
    "yourmodule/internal/service"
    "yourmodule/internal/controller"
    "yourmodule/internal/routes"
)

func main() {
    // 1. 连接数据库
    dsn := "root:password@tcp(localhost:3306)/testdb"
    db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    // 2. 初始化依赖
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo, db, nil)
    userController := controller.NewUserController(userService)

    // 3. 设置路由
    r := gin.Default()
    controllers := &routes.Controllers{
        User: userController,
    }
    routes.RegisterAutoRoutes(r, controllers)

    // 4. 启动服务
    r.Run(":8080")
}
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
