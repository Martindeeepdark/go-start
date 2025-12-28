# go-start 更新指南

## 🔄 如何更新到最新版本

### 方式 1: 一键更新脚本 (推荐)

```bash
curl -sSL https://raw.githubusercontent.com/Martindeeepdark/go-start/main/update.sh | bash
```

这个脚本会:
- ✅ 显示当前版本
- ✅ 自动下载并安装最新版本
- ✅ 处理 Go 1.24+ 的平台子目录
- ✅ 创建符号链接
- ✅ 验证安装成功

### 方式 2: 使用 go install

```bash
# 使用直连方式避免代理缓存
GOPROXY=direct go install github.com/Martindeeepdark/go-start/cmd/go-start@latest

# 验证版本
go-start --version
```

### 方式 3: 从源码编译

```bash
# 克隆仓库
git clone https://github.com/Martindeeepdark/go-start.git
cd go-start

# 编译
go build -o go-start ./cmd/go-start

# 安装
sudo mv go-start /usr/local/bin/
# 或者
mv go-start $GOPATH/bin/
```

---

## 📋 更新前后对比

### 查看当前版本

```bash
$ go-start --version
go-start 版本 v0.1.0
```

### 更新后

```bash
$ curl -sSL https://raw.githubusercontent.com/Martindeeepdark/go-start/main/update.sh | bash
📌 当前版本: v0.1.0
⬇️  正在更新 go-start...
✅ 安装命令执行成功
✅ go-start 已更新到: /Users/wenyz/go/bin/darwin_arm64/go-start
📌 新版本: v0.1.1
🎉 更新完成!
```

### 验证更新

```bash
$ go-start --version
go-start 版本 v0.1.1
```

---

## 🔧 常见问题

### Q1: 更新后还是旧版本怎么办?

**原因:** Go 模块缓存问题

**解决:**
```bash
# 清理 Go 模块缓存
go clean -modcache

# 重新安装
GOPROXY=direct go install github.com/Martindeeepdark/go-start/cmd/go-start@latest

# 验证
go-start --version
```

### Q2: 提示命令不存在?

**原因:** `$GOPATH/bin` 不在 PATH 中

**解决:**
```bash
# 临时解决
export PATH=$PATH:$(go env GOPATH)/bin

# 永久解决 (添加到 ~/.zshrc 或 ~/.bashrc)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

### Q3: 网络问题导致下载失败?

**原因:** 无法访问 GitHub 或 Go proxy 缓存问题

**解决:**
```bash
# 使用直连方式
GOPROXY=direct go install github.com/Martindeeepdark/go-start/cmd/go-start@latest

# 或者设置国内代理
export GOPROXY=https://goproxy.cn,direct
go install github.com/Martindeeepdark/go-start/cmd/go-start@latest
```

### Q4: Go 1.24+ 用户找不到可执行文件?

**原因:** Go 1.24+ 将可执行文件安装到平台子目录

**解决:**
```bash
# 查找实际安装位置
ls $(go env GOPATH)/bin/*/go-start

# 创建符号链接
ln -s $(go env GOPATH)/bin/darwin_arm64/go-start $(go env GOPATH)/bin/go-start

# 验证
go-start --version
```

---

## 📊 版本历史

| 版本 | 发布日期 | 主要变更 |
|------|---------|---------|
| **v0.1.1** | 2025-12-28 | 🐛 修复 create 命令 pkg 文件复制错误 |
| **v0.1.0** | 2025-12-28 | 🎯 全中文 CLI,增强 doctor,改进配置 |
| **v0.0.7** | 2025-12-28 | ✨ 简化向导流程 |

---

## 🔄 自动更新 (可选)

如果你想要类似 `brew upgrade` 的体验,可以创建一个 alias:

```bash
# 添加到 ~/.zshrc 或 ~/.bashrc
alias go-start-upgrade='curl -sSL https://raw.githubusercontent.com/Martindeeepdark/go-start/main/update.sh | bash'
```

然后就可以直接运行:
```bash
go-start-upgrade
```

---

## 💡 最佳实践

1. **定期检查更新**
   ```bash
   # 查看当前版本
   go-start --version

   # 查看 GitHub 最新版本
   curl -s https://api.github.com/repos/Martindeeepdark/go-start/releases/latest | grep '"tag_name"'
   ```

2. **查看更新日志**
   ```bash
   # 访问 Release 页面
   open https://github.com/Martindeeepdark/go-start/releases
   ```

3. **测试新版本**
   ```bash
   # 创建测试项目
   cd /tmp
   go-start create test-update

   # 验证功能正常
   cd test-update
   ls -la
   ```

---

## 🆘 获取帮助

如果更新过程中遇到问题:

- 📖 **文档**: [README.md](https://github.com/Martindeeepdark/go-start#readme)
- 🐛 **问题反馈**: [GitHub Issues](https://github.com/Martindeeepdark/go-start/issues)
- 💬 **讨论**: [GitHub Discussions](https://github.com/Martindeeepdark/go-start/discussions)

---

**更新愉快!** 🚀
