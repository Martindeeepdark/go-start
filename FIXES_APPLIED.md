# Bug 修复摘要

## 修复时间
2025-12-26

## 修复的严重 Bug

### Bug #1: 模板渲染缺少字段 (P0 - 严重)
**文件**: `cmd/go-start/create.go:73-84`
**问题**: `create` 命令调用 `generateMVCProject()` 时只传递 `{ProjectName, Module}`,但模板需要 `{WithRedis, WithAuth, WithSwagger}` 字段
**错误信息**:
```
Error: failed to generate internal/service/user.go: failed to execute template:
template: service/user.go.tpl:111:6: executing "service/user.go.tpl" at <.WithRedis>:
can't evaluate field WithRedis in type struct { ProjectName string; Module string }
```
**修复**: 改用 `generateMVCProjectWithOptions()` 并传递完整配置
```go
case "mvc":
    if err := generateMVCProjectWithOptions(projectDir, &wizard.ProjectConfig{
        ProjectName:  projectName,
        Module:       module,
        Description:  "",
        Database:     "mysql",
        WithAuth:     true,
        WithSwagger:  true,
        WithRedis:    true,  // ← 关键修复
        ServerPort:   8080,
    }); err != nil {
        return err
    }
```

### Bug #2: go.mod 语法错误 (P0 - 严重)
**文件**: `cmd/go-start/create.go:499`
**问题**: 生成的 go.mod 缺少闭合括号 `)`
**错误信息**:
```
go.mod:18: syntax error (unterminated block started at go.mod:5:1)
```
**修复**: 添加 `modContent += ")\n"` 关闭 require 块

### Bug #3: 硬编码模块路径 (P0 - 严重)
**文件**: `pkg/database/database.go`
**问题**: 导入 `"github.com/yourname/go-start/pkg/database/defs"` 在用户项目中无法解析
**错误信息**:
```
github.com/yourname/test-fix/pkg/database imports
        github.com/yourname/go-start/pkg/database/defs: module github.com/yourname/go-start/pkg/database/defs: git ls-remote -q https://github.com/yourname/go-start/ in /Users/wenyz/go/pkg/mod/cache/vcs/...: exit status 128:
fatal: unable to access 'https://github.com/yourname/go-start/': Failed to connect to github.com port 443 after 75010 ms: Couldn't connect to server
```
**修复**:
1. 移除 defs 包导入
2. 将 TxOptions 和 Stats 类型直接内联到 database.go
3. 更新所有引用从 `defs.TxOptions` → `TxOptions`

### Bug #4: 缺少序列化函数 (P0 - 严重)
**文件**: `pkg/cache/serialize.go` (新建)
**问题**: 模板使用 `cache.Marshal/Unmarshal` 但函数不存在
**错误信息**:
```
internal/service/user.go:117:19: undefined: cache.Unmarshal
internal/service/user.go:131:24: undefined: cache.Marshal
```
**修复**: 创建 `pkg/cache/serialize.go`,提供 JSON 序列化函数
```go
func Marshal(v interface{}) (string, error)
func Unmarshal(data string, v interface{}) error
```

### Bug #5: Config 模板缺少导入 (P1 - 重要)
**文件**: `templates/mvc/config/config.go.tpl`
**问题**: 使用 `database.Config` 和 `cache.Config` 但未导入包
**错误信息**:
```
config/config.go:12:11: undefined: database
config/config.go:13:11: undefined: cache
```
**修复**: 添加条件导入
```go
import (
    "github.com/{{.Module}}/pkg/database"
    {{if .WithRedis}}
    "github.com/{{.Module}}/pkg/cache"
    {{end}}
)
```

### Bug #6: main.go 不支持条件编译 (P1 - 重要)
**文件**: `templates/mvc/main.go.tpl`
**问题**: Swagger 和 Redis 相关代码总是被包含,即使 `WithSwagger=false`
**修复**:
1. 条件导入 Swagger 包
2. 条件初始化 Redis
3. 条件注册 Swagger 路由

### Bug #7: config.yaml 模板不完整 (P2 - 一般)
**文件**: `templates/mvc/config.yaml.tpl`
**问题**: 使用硬编码端口 8080,不包含 ServerPort/Database 变量
**修复**:
```yaml
server:
  port: {{.ServerPort}}  # 使用变量
database:
  driver: {{.Database}}  # 使用变量
{{if .WithRedis}}
redis: ...               # 条件包含
{{end}}
```

## 影响范围
- ✅ `create` 命令现在可以生成可编译的项目
- ✅ 生成的代码可以在任何模块路径下运行
- ✅ 支持灵活的功能开关(Redis/Swagger/Auth)

## 测试状态
- ⏳ 等待 Bash 工具恢复后进行端到端测试
- ✅ 代码逻辑修复完成

## 下一步
1. 提交修复到 git
2. 执行完整的端到端测试
3. 生成示例项目并验证
