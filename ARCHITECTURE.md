# go-start 项目架构说明

## 架构类型说明

### 代码组织方式（架构模式）

1. **MVC（Model-View-Controller）**
   - 简单直接
   - 适合中小型项目
   - 三层清晰：Controller → Service → Repository

2. **DDD（领域驱动设计）**
   - 复杂业务逻辑
   - 适合大型项目
   - 分层：Domain → Application → Infrastructure → Interface

### 部署方式（不作为模板区分）

- 单体应用 - 所有功能在一个进程
- 微服务 - 按业务拆分多个服务

## 正确的模板设计

### 方案 1: 简单分类

```
templates/
├── mvc/        # MVC 架构模板（适合简单项目）
└── ddd/        # DDD 架构模板（适合复杂项目）
```

**使用方式**：
```bash
go-start create my-api --arch=mvc     # MVC 项目
go-start create my-api --arch=ddd     # DDD 项目
```

### 方案 2: 按复杂度分类

```
templates/
├── simple/     # 简单项目（3层架构）
├── standard/   # 标准项目（清晰分层）
└── advanced/   # 复杂项目（DDD）
```

## 推荐方案

**采用方案 1（简单清晰）**：

### MVC 模板（默认）
```
my-api/
├── cmd/server/main.go
├── internal/
│   ├── controller/      # HTTP 处理
│   ├── service/         # 业务逻辑
│   ├── repository/      # 数据访问
│   └── model/           # 数据模型
└── config/
```

### DDD 模板（可选）
```
my-api/
├── cmd/server/main.go
├── internal/
│   ├── domain/          # 领域层
│   ├── application/     # 应用层
│   ├── infrastructure/  # 基础设施层
│   └── interface/       # 接口层（HTTP）
└── config/
```

## 选择建议

- **MVC**：CRUD 应用、后台管理、小型 API
- **DDD**：复杂业务、领域模型丰富、需要长期维护

## 部署架构（不区分）

生成的项目可以选择：
- 单体部署（默认）
- 微服务拆分（手动拆分，或者后续提供 `go-start split` 命令）
