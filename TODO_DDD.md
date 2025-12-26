# DDD 架构实现 TODO

## 当前状态
- ✅ DDD 模板已创建 (`templates/ddd/`)
- ✅ DDD 代码生成器已实现 (`pkg/gen/ddd.go`)
- ✅ `gen db --arch=ddd` 命令支持 DDD
- ❌ `create --arch=ddd` 命令未实现

## 需要实现的功能

### 1. 实现 generateDDDProjectWithOptions()
**文件**: `cmd/go-start/create.go`

参考 `generateMVCProjectWithOptions()`,创建:
```go
func generateDDDProjectWithOptions(projectDir string, config *wizard.ProjectConfig) error {
    // 1. 创建 DDD 目录结构
    dirs := []string{
        "cmd/server",
        "internal/domain",
        "internal/application",
        "internal/infrastructure/persistence",
        "internal/interface/http",
        "config",
        "pkg/...",
    }

    // 2. 使用 DDD 模板生成文件
    templateFiles := map[string]string{
        "cmd/server/main.go": "main.go.tpl",
        // ... 其他 DDD 模板
    }

    // 3. 复制 pkg 文件
    copyPkgFiles(projectDir)

    return nil
}
```

### 2. 更新 create.go 的 switch 语句
**文件**: `cmd/go-start/create.go:122`

```go
case "ddd":
    if err := generateDDDProjectWithOptions(projectDir, config); err != nil {
        return err
    }
```

### 3. DDD 模板完善
检查以下模板是否完整:
- `templates/ddd/main.go.tpl` ✅ (已创建,但内容可能不完整)
- `templates/ddd/README.md.tpl` ✅
- `templates/ddd/internal/domain/*/` ✅
- `templates/ddd/internal/application/*/` ✅
- `templates/ddd/internal/infrastructure/*/` ✅
- `templates/ddd/internal/interface/*/` ✅

### 4. DDD 项目示例
创建一个完整的 DDD 示例项目,验证模板正确性。

## 优先级
P1 - 高优先级,但次于当前 bug 修复和测试

## 参考资料
- `pkg/gen/ddd.go` - DDD 代码生成器实现
- `templates/ddd/` - DDD 项目模板
- `docs/DDD_GUIDE.md` - DDD 使用指南
