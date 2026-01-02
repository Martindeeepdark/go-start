# GitHub Actions 工作流

本项目使用 GitHub Actions 进行自动化构建、测试和发布。

## 工作流列表

### 1. Go 工作流 (`.github/workflows/go.yml`)

完整的 CI/CD 流程，包括：

#### 代码质量检查 ✅
- **依赖验证**: `go mod verify`
- **代码检查**: `go vet`
- **格式化检查**: `gofmt`
- **单元测试**: `go test -race -coverprofile`
- **代码覆盖率**: 自动上传到 Codecov

#### 跨平台编译 🔨
支持以下平台的二进制文件构建：
- **Linux**: AMD64, ARM64
- **macOS**: AMD64, ARM64 (Apple Silicon)
- **Windows**: AMD64

#### 自动发布 🚀
当推送带 `v` 前缀的标签时（如 `v1.2.0`），自动：
1. 编译所有平台的二进制文件
2. 生成 Release Notes
3. 创建 GitHub Release
4. 上传所有构建产物

### 2. CodeQL 工作流 (`.github/workflows/codeql.yml`)

代码安全分析，用于发现潜在的安全漏洞和代码质量问题。

## 触发条件

### 自动触发
- **推送到 main 分支**: 运行测试和构建
- **Pull Request**: 运行测试
- **推送标签 (v*)**: 运行完整流程并发布 Release
- **每周日**: 运行 CodeQL 分析

### 手动触发
可以在 GitHub Actions 页面手动运行工作流。

## 使用方法

### 正常开发
```bash
# 1. 开发并提交代码
git add .
git commit -m "feat: 添加新功能"
git push origin main

# → GitHub Actions 自动运行测试
```

### 发布新版本
```bash
# 1. 更新版本号（可选）
# 编辑 cmd/go-start/main.go 中的 Version 变量

# 2. 提交更改
git add .
git commit -m "chore: 更新版本号到 v1.3.0"
git push origin main

# 3. 创建并推送标签
git tag -a v1.3.0 -m "v1.3.0: 发布说明"
git push origin v1.3.0

# → GitHub Actions 自动：
#   - 运行所有测试
#   - 编译所有平台的二进制文件
#   - 创建 GitHub Release
#   - 上传构建产物
```

## 构建产物

每次构建生成的二进制文件命名规则：

```
go-start-{OS}-{ARCH}.{ext}
```

例如：
- `go-start-linux-amd64.tar.gz`
- `go-start-darwin-arm64.tar.gz` (Apple Silicon)
- `go-start-windows-amd64.zip`

## 下载安装

从 Release 页面下载对应平台的文件：

```bash
# Linux / macOS
tar xzf go-start-linux-amd64.tar.gz
sudo mv go-start /usr/local/bin/
sudo chmod +x /usr/local/bin/go-start

# Windows
# 解压 go-start-windows-amd64.zip
# 将 go-start.exe 移动到 PATH 目录
```

## 状态徽章

在 README.md 中添加状态徽章：

```markdown
![Build Status](https://github.com/Martindeeepdark/go-start/actions/workflows/go.yml/badge.svg)
![CodeQL](https://github.com/Martindeeepdark/go-start/actions/workflows/codeql.yml/badge.svg)
```

## 环境变量

工作流中使用的环境变量：

- `GO_VERSION`: Go 版本 (默认: 1.21)
- `CGO_ENABLED`: 0 (静态编译)
- `GITHUB_TOKEN`: 自动提供，用于创建 Release

## 故障排查

### 构建失败
1. 检查日志: Actions 页面 → 选择工作流 → 查看详细日志
2. 本地测试: `go test ./...` 和 `go build ./...`
3. 检查格式: `gofmt -l .`

### Release 创建失败
1. 确认标签格式正确: `v*.*.*`
2. 确认 GITHUB_TOKEN 权限足够
3. 检查构建产物是否成功生成

## 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Go 在 GitHub Actions 中的最佳实践](https://github.com/actions/setup-go)
- [CodeQL 文档](https://codeql.github.com/docs/)
