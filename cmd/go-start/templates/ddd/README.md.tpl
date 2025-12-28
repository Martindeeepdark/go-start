# {{.ProjectName}}

{{if .Description}}{{.Description}}{{else}}{{.ProjectName}} - 基于 DDD 架构的 Go Web 服务{{end}}

## 架构说明

本项目采用 **DDD（领域驱动设计）** 架构，分为四层：

### 1. Domain 层（领域层）
- 封装业务逻辑和规则
- 包含：实体（Entity）、仓储接口（Repository）、领域服务（Domain Service）

### 2. Application 层（应用层）
- 编排用例，协调领域对象
- 包含：应用服务（Application Service）

### 3. Infrastructure 层（基础设施层）
- 提供技术实现
- 包含：仓储实现（Repository Implementation）

### 4. Interface 层（接口层）
- 对外接口，处理协议
- 包含：HTTP 控制器（Controller）、路由

## 快速开始

```bash
# 安装依赖
go mod tidy

# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 编辑配置
vim config/config.yaml

# 运行服务
go run cmd/server/main.go
```

## 项目结构

```
.
├── cmd/
│   └── server/
│       └── main.go           # 应用入口
├── internal/
│   ├── domain/               # 领域层
│   │   ├── user/             # 用户聚合
│   │   │   ├── User.go       # 用户实体
│   │   │   ├── repository.go # 仓储接口
│   │   │   └── service.go    # 领域服务
│   ├── application/          # 应用层
│   │   └── user/             # 用户应用服务
│   ├── infrastructure/       # 基础设施层
│   │   └── persistence/      # 持久化
│   │       └── user/         # 用户仓储实现
│   └── interface/            # 接口层
│       └── http/             # HTTP 接口
│           ├── user/         # 用户控制器
│           └── routes.go     # 路由
├── config/                   # 配置
├── pkg/                      # 工具包
└── scripts/                  # 脚本
```

## 开发指南

### 添加新功能

1. **Domain 层**：定义实体和业务规则
2. **Application 层**：编排用例
3. **Infrastructure 层**：实现持久化
4. **Interface 层**：暴露 HTTP 接口

### 代码生成

使用 go-start 从数据库生成代码：

```bash
# 生成 DDD 架构代码
go-start gen db --dsn="..." --tables=users --arch=ddd
```

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL / PostgreSQL
- **缓存**: Redis
- **日志**: Zap
- **配置**: Viper

## 许可证

MIT License
