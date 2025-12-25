# go-start 架构设计

## 1. 项目模板设计

### 1.1 创建命令选项

```bash
go-start create my-api --arch=mvc      # MVC 架构
go-start create my-api --arch=ddd      # DDD 架构
go-start create my-api --arch=msa      # 微服务架构
```

### 1.2 MVC 模板结构

```
my-api/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── controller/              # 控制器层
│   │   ├── user_controller.go
│   │   └── base.go
│   ├── service/                 # 业务逻辑层
│   │   ├── user_service.go
│   │   └── base.go
│   ├── repository/              # 数据访问层
│   │   ├── user_repository.go
│   │   └── base.go
│   ├── model/                   # 数据模型
│   │   ├── user.go
│   │   └── base.go
│   └── middleware/              # 中间件
│       ├── auth.go
│       └── logger.go
├── config/                      # 配置
│   ├── config.yaml
│   └── config.go
├── pkg/                         # 公共包
│   └── response/                # 统一响应
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 1.3 DDD 模板结构

```
my-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/                  # 领域层
│   │   ├── user/               # 用户领域
│   │   │   ├── entity.go       # 实体
│   │   │   ├── valueobject.go  # 值对象
│   │   │   ├── repository.go   # 仓储接口
│   │   │   └── service.go      # 领域服务
│   │   └── order/              # 订单领域
│   ├── application/             # 应用层
│   │   ├── service/            # 应用服务
│   │   ├── dto/                # 数据传输对象
│   │   └── assembler/          # 组装器
│   ├── infrastructure/          # 基础设施层
│   │   ├── persistence/        # 持久化实现
│   │   │   ├── mysql/
│   │   │   └── redis/
│   │   ├── mq/                 # 消息队列
│   │   └── cache/              # 缓存实现
│   └── interface/              # 接口层（HTTP/gRPC）
│       ├── handler/            # 处理器
│       └── middleware/         # 中间件
├── config/
├── pkg/
└── README.md
```

### 1.4 微服务模板 (MSA)

```
my-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/                  # 领域层（轻量）
│   ├── service/                 # 服务层
│   ├── repository/              # 仓储层
│   ├── interface/              # 接口层
│   └── middleware/
├── config/
├── pkg/
│   ├── httpx/                  # HTTP 工具
│   ├── database/               # 数据库工具
│   └── cache/                  # 缓存工具
└── README.md
```

---

## 2. 目录结构设计

### 2.1 go-start 自身结构

```
go-start/
├── cmd/
│   └── go-start/
│       └── main.go               # CLI 入口
├── internal/
│   ├── cli/                     # CLI 命令
│   │   ├── create.go            # 创建项目命令
│   │   ├── run.go               # 运行命令
│   │   └── build.go             # 构建命令
│   ├── generator/               # 代码生成器
│   │   ├── project.go           # 项目生成
│   │   ├── api.go               # API 生成
│   │   └── model.go             # Model 生成
│   └── config/                  # 脚手架配置
├── pkg/                         # 给生成项目的包
│   ├── httpx/                   # Gin 工具（从 common 移过来）
│   ├── database/                # GORM 集成
│   │   ├── mysql.go             # MySQL 实现
│   │   ├── postgres.go          # PostgreSQL 实现
│   │   └── db.go                # 数据库管理
│   ├── cache/                   # Redis 集成
│   │   ├── redis.go             # Redis 客户端
│   │   ├── hash.go              # Hash 操作
│   │   ├── list.go              # List 操作
│   │   └── lock.go              # 分布式锁
│   └── middleware/              # 中间件
│       ├── auth.go
│       ├── logger.go
│       └── recovery.go
├── templates/                  # 项目模板
│   ├── mvc/                    # MVC 模板
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── config/
│   │   └── template.go
│   ├── ddd/                    # DDD 模板
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── config/
│   │   └── template.go
│   └── msa/                    # 微服务模板
│       ├── cmd/
│       ├── internal/
│       ├── config/
│       └── template.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 3. 核心模块设计

### 3.1 httpx 模块（从 common 移过来并增强）

```go
// pkg/httpx/response/response.go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{})
func Error(c *gin.Context, code int, message string)
func Paginated(c *gin.Context, items interface{}, total int64, page, pageSize int)
```

### 3.2 database 模块

```go
// pkg/database/db.go
type DB struct {
    *gorm.DB
}

func InitMySQL(cfg *config.MySQLConfig) (*DB, error)
func InitPostgres(cfg *config.PostgresConfig) (*DB, error)

// 实现 common/database/defs 接口
func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (defs.Rows, error)
func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (defs.Result, error)
func (db *DB) BeginTx(ctx context.Context, opts *defs.TxOptions) (defs.Transaction, error)
```

### 3.3 cache 模块

```go
// pkg/cache/redis.go
type RedisCache struct {
    client *redis.Client
}

func NewRedis(cfg *config.RedisConfig) (*RedisCache, error)

// 实现 common/cache/defs 接口
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
// ... 其他接口实现
```

---

## 4. 配置设计

### 4.1 配置文件 (config/config.yaml)

```yaml
server:
  port: 8080
  mode: debug

database:
  driver: mysql
  mysql:
    host: localhost
    port: 3306
    username: root
    password: ""
    database: my_api
    max_open_conns: 100
    max_idle_conns: 10

  postgres:
    host: localhost
    port: 5432
    username: postgres
    password: ""
    database: my_api

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

log:
  level: info
  output: stdout
```

### 4.2 配置加载

```go
// config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Log      LogConfig
}

func Load(path string) (*Config, error)
```

---

## 5. 使用流程

### 5.1 创建新项目

```bash
# MVC 架构
go-start create my-api --arch=mvc

# DDD 架构
go-start create my-api --arch=ddd

# 微服务架构（轻量 DDD）
go-start create my-api --arch=msa
```

### 5.2 生成代码

```bash
cd my-api

# 生成 API
go-start api add user

# 生成 Model
go-start model add User

# 运行
go-start run
```

---

## 6. 优先级

### Phase 1 (MVP - 本周完成)
- [ ] go-start 项目初始化
- [ ] MVC 模板
- [ ] httpx 模块
- [ ] database 模块（MySQL）
- [ ] cache 模块（Redis）
- [ ] 基础 CLI 命令（create/run）

### Phase 2 (增强 - 下周)
- [ ] DDD 模板
- [ ] MSA 模板
- [ ] 代码生成器（api/model）
- [ ] 配置增强

### Phase 3 (高级 - 后续)
- [ ] eventbus 集成
- [ ] Swagger 生成
- [ ] Docker 生成
- [ ] 监控集成
