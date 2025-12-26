# go-start 项目结构

```
go-start/
├── cmd/
│   └── go-start/              # CLI 应用入口
│       ├── main.go            # 主程序
│       ├── create.go          # 创建项目命令
│       ├── run.go             # 运行项目命令
│       ├── spec.go            # spec-kit 规范生成命令
│       └── gen.go             # 数据库 CRUD 生成命令
│
├── pkg/                       # 核心功能包
│   ├── wizard/                # 交互式向导
│   ├── spec/                  # spec-kit 规范解析和代码生成
│   ├── gen/                   # 数据库 CRUD 代码生成
│   ├── cache/                 # Redis 缓存封装（模板用）
│   ├── database/              # 数据库封装（模板用）
│   ├── httpx/                 # HTTP 工具（模板用）
│   └── middleware/            # 中间件（模板用）
│
├── templates/                 # 项目模板
│   ├── mvc/                   # MVC 架构模板
│   └── ddd/                   # DDD 架构模板（预留）
│
├── spec/                      # spec-kit 示例规范文件
│   └── example.blog.spec.yaml
│
├── docs/                      # 文档（建议创建）
│   ├── README.md
│   ├── WIZARD.md
│   ├── SPEC.md
│   └── ARCHITECTURE.md
│
├── go.mod                     # Go 模块定义
├── go.sum
├── Makefile                   # 构建脚本
├── .gitignore
└── README.md
```

## 目录说明

### cmd/
CLI 应用的入口点，所有命令行工具的代码。

### pkg/
核心功能代码库：
- **wizard/** - 交互式项目创建向导
- **spec/** - YAML 规范解析和代码生成
- **gen/** - 从数据库生成 CRUD 代码
- 其他包是项目生成时复制到新项目的模板代码

### templates/
项目模板文件，用于生成新项目的代码结构。

### spec/
示例规范文件，展示如何使用 spec-kit。

## 清理规则

- ✅ **保留** - 核心代码和文档
- ❌ **删除** - 测试生成的项目（*-api/）
- ❌ **删除** - 编译产物（bin/）
- ❌ **删除** - 临时文件
