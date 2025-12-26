# GORM Gen 代码生成器使用指南

## 功能说明

go-start 集成了 GORM Gen（生产环境验证的 GORM 代码生成器），可以自动从数据库表生成完整的 CRUD 代码。

## 使用方法

### 1. 交互式选择表（推荐）

```bash
go-start gen db --dsn="root:password@tcp(localhost:3306)/mydb" --interactive
```

这个命令会：
1. 连接数据库
2. 显示所有表及其字段数、索引数
3. 让你交互式选择要生成的表

交互式选择支持多种输入方式：
- 数字序号：`1,2,3` - 选择第 1、2、3 张表
- 范围：`1-5` - 选择第 1 到第 5 张表
- 通配符：`user*` - 选择所有以 user 开头的表
- 全部：`all` - 选择所有表

### 2. 直接指定表名

```bash
go-start gen db --dsn="root:password@tcp(localhost:3306)/mydb" --tables=users,articles,comments
```

### 3. 使用通配符

```bash
# 生成所有以 user 开头的表
go-start gen db --dsn="..." --tables="user*"

# 生成所有包含 _log 的表
go-start gen db --dsn="..." --tables="*_log"
```

### 4. 自定义输出目录

```bash
go-start gen db --dsn="..." --tables=users --output=./myproject/internal
```

## 生成的代码结构

执行命令后，GORM Gen 会生成以下代码：

```
internal/
└── dal/
    ├── query/
    │   ├── gen.go           # GORM Gen 生成的查询 API
    │   ├── users.go         # User 表的查询方法
    │   ├── articles.go      # Article 表的查询方法
    │   └── ...
    └── models.go            # 数据模型定义
```

### GORM Gen 生成的查询 API

每个表都会生成对应的查询 API，包括：

```go
// 基础查询
query.User.WithContext(ctx).Where(&User{Name: "jinzhu"}).First()

// 条件查询
query.User.WithContext(ctx).Where(query.User.Name.Eq("jinzhu")).First()

// 动态条件
u := query.User
u.WithContext(ctx).Where(u.Age.Gt(18)).Find()

// 关联查询
u.WithContext(ctx).Preload(u.Orders).First()
```

## 支持的数据库

### MySQL

```bash
go-start gen db --dsn="user:password@tcp(host:3306)/dbname" --tables=users
```

### PostgreSQL

```bash
go-start gen db --dsn="host=localhost user=postgres password=postgres dbname=testdb port=5432" --tables=users
```

## GORM Gen 的优势

1. **生产验证**: GORM Gen 是 GORM 官方的代码生成器，经过大量生产环境验证
2. **类型安全**: 生成的代码是类型安全的，编译时检查
3. **API 友好**: 提供链式查询 API，代码更优雅
4. **性能优化**: 生成的代码性能优秀，无反射开销
5. **自动维护**: 数据库变更后重新生成即可

## 示例：从数据库生成完整 CRUD

```bash
# 1. 假设你有一个 MySQL 数据库
mysql -u root -p
CREATE DATABASE testdb;
USE testdb;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    UNIQUE KEY uk_email (email)
);

# 2. 使用 go-start 生成代码
go-start gen db --dsn="root:pass@tcp(localhost:3306)/testdb" --tables=users

# 3. 查看生成的代码
tree internal/dal

# 4. 在项目中使用生成的代码
# import "yourproject/internal/dal/query"
# import "yourproject/internal/dal/model"
```

## 下一步

当前版本仅生成 GORM Gen 基础代码（Model + Query API）。

后续将实现：
- Repository 层封装（基于 GORM Gen API）
- Service 层生成（包含业务逻辑 + 缓存）
- Controller 层生成（RESTful API 处理器）
- Routes 自动注册

这样你就可以完全专注于业务逻辑，无需手写重复的 CRUD 代码。
