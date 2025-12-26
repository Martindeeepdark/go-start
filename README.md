# go-start - 高级 Go 脚手架工具

> 帮助新人工程师快速上手，让高级工程师大展身手

## 项目愿景

开发一个高级版的 Go 脚手架工具（基于 nunu 的改进），目标是：

- **新人友好**: 详细中文注释、交互式向导、自动化工具
- **高手赋能**: spec-driven 开发、插件系统、完全可定制
- **CRUD 自动化**: 使用生产级工具自动生成代码，专注业务逻辑
- **本地运行**: 无需 GitHub/GitLab，完全本地化能力

## 核心功能

### 1. 交互式项目创建向导

```bash
go-start create
```

**特性：**
- 9 步向导式项目创建
- 全中文界面
- 支持 MVC/DDD 架构选择
- 支持 MySQL/PostgreSQL/SQLite
- 可选 Redis/认证/Swagger

### 2. GORM Gen 代码生成

```bash
# 交互式选择表
go-start gen db --dsn="root:pass@tcp(localhost:3306)/mydb" --interactive

# 指定表名（MVC 架构）
go-start gen db --dsn="..." --tables=users,articles --arch=mvc

# 使用 DDD 架构
go-start gen db --dsn="..." --tables=users --arch=ddd

# 通配符匹配
go-start gen db --dsn="..." --tables="user*"
```

**生成内容：**

#### MVC 架构（默认）
- ✅ GORM Gen Model 和 Query API
- ✅ Repository 层（基于 GORM Gen API，带中文注释）
- ✅ 基于索引自动生成查询方法
- ✅ Service 层（业务逻辑 + 缓存支持）
- ✅ Controller 层（RESTful API）
- ✅ 路由自动注册

#### DDD 架构
- ✅ Domain 层（实体、仓储接口、领域服务）
- ✅ Application 层（应用服务、用例编排）
- ✅ Infrastructure 层（仓储实现、持久化）
- ✅ Interface 层（HTTP 控制器、路由注册）

详见：[DDD 架构指南](./docs/DDD_GUIDE.md)

### 3. Spec-Kit 规范驱动开发

```bash
go-start spec generate --file=api.spec.yaml --output=./internal
```

**特性：**
- YAML 定义 API 规范
- 自动生成完整代码
- 支持自定义模板

## 技术栈

### 核心依赖
- **CLI**: [Cobra](https://github.com/spf13/cobra) - 命令行框架
- **Web**: [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- **ORM**: [GORM](https://github.com/go-gorm/gorm) + [GORM Gen](https://github.com/go-gorm/gen) - ORM 和代码生成
- **缓存**: [go-redis](https://github.com/redis/go-redis) - Redis 客户端

### 为什么选择 GORM Gen？

虽然 GORM 官方推出了新的 [GORM CLI](https://github.com/go-gorm/cli)，但我们选择使用成熟的 GORM Gen：

- ✅ **生产验证**: 经过大量生产环境验证
- ✅ **稳定可靠**: 版本 v0.3.27，功能完整
- ✅ **社区支持**: 文档完善，问题容易解决
- ✅ **类型安全**: 编译时检查，无运行时错误

详细对比请查看：[GORM 技术选型文档](./docs/GORM_TECH_CHOICE.md)

## 快速开始

### 安装

```bash
git clone https://github.com/yourname/go-start.git
cd go-start
go build -o bin/go-start cmd/go-start/*.go
```

### 创建项目

```bash
# 交互式创建（推荐）
./bin/go-start create

# 或指定配置
./bin/go-start create --name=my-api --arch=mvc --db=mysql
```

### 生成代码

```bash
# 1. 创建数据库表
mysql -u root -p
CREATE DATABASE testdb;
USE testdb;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_username (username)
);

# 2. 生成代码
./bin/go-start gen db \
    --dsn="root:pass@tcp(localhost:3306)/testdb" \
    --tables=users \
    --output=./internal

# 3. 查看生成的代码
tree internal/dal
tree internal/repository
```

## 项目结构

```
go-start/
├── cmd/go-start/          # CLI 工具
│   ├── main.go            # 入口
│   ├── create.go          # create 命令
│   ├── gen.go             # gen 命令（数据库生成）
│   └── spec.go            # spec 命令（规范驱动）
│
├── pkg/
│   ├── wizard/            # 交互式向导
│   ├── spec/              # spec-kit 规范解析
│   ├── gen/               # 数据库代码生成
│   ├── cache/             # 缓存封装
│   ├── database/          # 数据库管理
│   └── httpx/             # HTTP 工具
│
├── templates/
│   └── mvc/               # MVC 架构模板
│       ├── model/
│       ├── repository/
│       ├── service/
│       └── controller/
│
├── spec/                  # 示例规范文件
└── docs/                  # 文档
    ├── GORM_TECH_CHOICE.md    # GORM 技术选型
    ├── GORM_GEN_GUIDE.md      # GORM Gen 使用指南
    └── PROJECT_STATUS.md      # 项目状态
```

## 生成的代码示例

### Repository 层

```go
// UserRepository 用户数据访问层
//
// 职责说明：
//   - 封装 User 的数据库操作
//   - 提供基础 CRUD 方法
//   - 基于索引生成高效查询方法
//   - 使用 GORM Gen 生成的类型安全 API
type UserRepository struct {
    q *query.Query
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    return r.q.User.WithContext(ctx).
        Where(r.q.User.ID.Eq(id)).  // ✅ 类型安全，无魔法字符串
        First()
}

// GetByEmail 根据邮箱获取用户（使用索引）
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    return r.q.User.WithContext(ctx).
        Where(r.q.User.Email.Eq(email)).
        First()
}
```

## 设计理念

### 新人友好
- **详细注释**: 每个函数都有中文注释说明职责、参数、返回值
- **自动化**: 减少重复劳动，专注业务逻辑
- **最佳实践**: 内置 Go 项目最佳实践
- **交互式向导**: 降低学习曲线

### 高手赋能
- **spec-driven**: YAML 定义 API，自动生成代码
- **可定制**: 支持自定义模板、插件
- **生产级**: 使用经过生产验证的工具（GORM Gen）
- **架构选择**: 支持 MVC 和 DDD

## 开发进度

### ✅ 已完成
- [x] 交互式项目创建向导
- [x] Spec-Kit 规范驱动开发
- [x] GORM Gen 集成
- [x] 数据库连接和表结构读取
- [x] Repository 层生成（带中文注释）
- [x] 基于索引自动生成查询方法
- [x] Service 层生成（业务逻辑 + 缓存）
- [x] Controller 层生成（RESTful API）
- [x] 路由自动注册
- [x] 完整使用示例和文档
- [x] DDD 架构支持
- [x] MVC/DDD 架构选择

### 🚧 进行中
- [ ] 优化生成代码的模块路径配置
- [ ] 添加更多单元测试

### 📋 待实现
- [ ] 认证系统（JWT）
- [ ] Swagger 文档生成
- [ ] 插件系统

详见：[项目状态文档](./docs/PROJECT_STATUS.md)

## 使用示例

### 完整工作流

```bash
# 1. 创建项目
./bin/go-start create
# 按向导填写配置

# 2. 设计数据库
# 在你的 MySQL 中创建表

# 3. 生成代码
./bin/go-start gen db --dsn="..." --tables=users --interactive

# 4. 编写业务逻辑
# 在生成的 Service 层中添加你的业务逻辑

# 5. 运行项目
go run cmd/server/main.go
```

## 文档

- [完整使用示例](./docs/COMPLETE_EXAMPLE.md) - 端到端的使用教程
- [DDD 架构指南](./docs/DDD_GUIDE.md) - DDD 架构使用指南
- [GORM 技术选型](./docs/GORM_TECH_CHOICE.md) - 为什么选择 GORM Gen
- [GORM Gen 使用指南](./docs/GORM_GEN_GUIDE.md) - 详细使用教程
- [项目状态](./docs/PROJECT_STATUS.md) - 开发进度和规划

## 贡献

欢迎贡献代码、提出建议或报告问题！

## 许可证

MIT License

---

**目标**: 让新人工程师快速上手，让高级工程师大展身手，让每个人都能专注于业务逻辑。
