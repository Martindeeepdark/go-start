#!/bin/bash

# go-start 更新脚本
# 用于更新已安装的 go-start 到最新版本

set -e

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                                                           ║"
echo "║   🔄 go-start 更新脚本                                    ║"
echo "║                                                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未检测到 Go 环境"
    echo ""
    echo "请先安装 Go:"
    echo "  macOS: brew install go"
    echo "  Linux: 下载 https://go.dev/dl/"
    echo ""
    exit 1
fi

# 显示当前版本
if command -v go-start &> /dev/null; then
    CURRENT_VERSION=$(go-start --version 2>&1 | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "未知")
    echo "📌 当前版本: $CURRENT_VERSION"
else
    echo "📌 go-start 未安装"
fi
echo ""

# 检测 Go bin 路径
GOPATH=$(go env GOPATH)
GOBIN="$GOPATH/bin"

# Go 1.24+ 会安装到平台特定子目录
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
PLATFORM_DIR="$GOBIN/${GOOS}_${GOARCH}"

echo "⬇️  正在更新 go-start..."
echo ""

# 使用直连方式避免代理缓存问题
if GOPROXY=direct go install github.com/Martindeeepdark/go-start/cmd/go-start@latest 2>&1; then
    echo "✅ 安装命令执行成功"
else
    echo "❌ 安装命令执行失败"
    echo ""
    echo "可能的原因:"
    echo "  1. 网络问题,无法访问 GitHub"
    echo "  2. Go 版本不兼容 (需要 Go 1.21+)"
    echo ""
    echo "请尝试手动安装:"
    echo "  git clone https://github.com/Martindeeepdark/go-start.git"
    echo "  cd go-start"
    echo "  go build -o go-start ./cmd/go-start"
    echo "  mv go-start /usr/local/bin/"
    exit 1
fi

# 检测实际安装位置
if [ -f "$PLATFORM_DIR/go-start" ]; then
    ACTUAL_BIN_PATH="$PLATFORM_DIR/go-start"
    # 创建符号链接 (如果不存在)
    if [ ! -L "$GOBIN/go-start" ] && [ ! -f "$GOBIN/go-start" ]; then
        ln -s "$PLATFORM_DIR/go-start" "$GOBIN/go-start"
        echo "🔗 创建符号链接: $GOBIN/go-start -> $PLATFORM_DIR/go-start"
    fi
elif [ -f "$GOBIN/go-start" ]; then
    ACTUAL_BIN_PATH="$GOBIN/go-start"
else
    echo "❌ 安装失败: 未找到可执行文件"
    echo "   检查了以下路径:"
    echo "   - $GOBIN/go-start"
    echo "   - $PLATFORM_DIR/go-start"
    exit 1
fi

echo ""
echo "✅ go-start 已更新到: $ACTUAL_BIN_PATH"
echo ""

# 显示新版本
NEW_VERSION=$($ACTUAL_BIN_PATH --version 2>&1 | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "未知")
echo "📌 新版本: $NEW_VERSION"
echo ""

# 检查是否在 PATH 中
if echo $PATH | grep -q "$GOBIN"; then
    echo "✅ $GOBIN 已在 PATH 中"
    echo ""
    echo "🎉 更新完成!可以使用:"
    echo "   go-start --version"
    echo "   go-start create my-api"
    echo ""
else
    echo "⚠️  $GOBIN 不在 PATH 中,请手动添加到 PATH"
    echo ""
    echo "📝 请运行以下命令:"
    echo "   export PATH=\"\$PATH:$GOBIN\""
    echo ""
    echo "   或者添加到你的 shell 配置文件中"
    echo ""
fi

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                                                           ║"
echo "║   ✅ 更新完成!                                            ║"
echo "║                                                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
