# 🐛 v0.1.1 - Bug 修复版本

这是 v0.1.0 的补丁版本,修复了一个重要的 bug。

## 🐛 修复的问题

### 严重 Bug: create 命令无法复制 pkg 文件

**问题描述:**
```bash
$ go-start create my-api
错误: 复制 pkg 文件失败: lstat /Users/wenyz/GolandProjects/pkg: no such file or directory
```

**根本原因:**
- 当从安装的二进制运行时,`getRootDir()` 返回错误的路径
- 导致无法找到 pkg 源码目录

**修复方案:**
- 添加 `findPkgDir()` 函数,智能查找 pkg 目录
- 支持从多种环境运行:
  1. ✅ 从源码运行 (`go run cmd/go-start/main.go`)
  2. ✅ 从本地二进制运行 (`./bin/go-start`)
  3. ✅ 从安装的二进制运行 (`go-start`)

**修复的代码:**
```go
func findPkgDir() string {
	// 1. 尝试当前目录 (开发时)
	if _, err := os.Stat("pkg"); err == nil {
		return filepath.Abs("pkg")
	}

	// 2. 尝试父目录 (从 cmd/go-star 运行时)
	parentPkg := filepath.Join("..", "..", "pkg")
	if _, err := os.Stat(parentPkg); err == nil {
		return filepath.Abs(parentPkg)
	}

	// 3. 尝试二进制的父目录 (从安装的二进制运行时)
	execDir := filepath.Dir(os.Args[0])
	paths := []string{
		filepath.Join(execDir, "..", "pkg"),
		filepath.Join(execDir, "..", "..", "pkg"),
		filepath.Join(execDir, "..", "..", "..", "pkg"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return filepath.Abs(path)
		}
	}

	return "" // 找不到时返回空,不阻塞创建流程
}
```

## 🧪 测试

### 测试场景 1: 从源码运行
```bash
$ go run cmd/go-start/main.go create test-project
✓ Project test-project created successfully!
```

### 测试场景 2: 从二进制运行
```bash
$ ./bin/go-start create test-project
✓ Project test-project created successfully!
```

### 测试场景 3: 从安装的二进制运行
```bash
$ go-start create test-project
✓ Project test-project created successfully!
```

## 📦 升级指南

### 从 v0.1.0 升级

```bash
# 方式 1: 一键安装脚本
curl -sSL https://raw.githubusercontent.com/Martindeeepdark/go-start/main/install.sh | bash

# 方式 2: 直接安装
GOPROXY=direct go install github.com/Martindeeepdark/go-start/cmd/go-start@latest

# 验证版本
go-start --version
# 应该显示: go-start 版本 v0.1.1
```

## 📝 完整变更

### 修复 (Bug Fixes)
- 修复 create 命令无法复制 pkg 文件的问题
- 改进 pkg 目录查找逻辑
- 添加多环境支持

## 🔗 相关链接

- **提交**: c43b3b9
- **标签**: v0.1.1
- **对比**: [v0.1.0...v0.1.1](https://github.com/Martindeeepdark/go-start/compare/v0.1.0...v0.1.1)

## 🙏 致谢

感谢反馈这个 bug!

---

**⭐ 如果这个项目对你有帮助,请给个 Star 支持一下!**

Made with ❤️ by [Martindeeepdark](https://github.com/Martindeeepdark)
