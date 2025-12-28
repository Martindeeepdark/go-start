#!/bin/bash

# go-start 一键安装脚本
# 适用于 macOS 和 Linux

set -e

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                                                           ║"
echo "║   🚀 go-start 一键安装脚本                                ║"
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

echo "✅ 检测到 Go: $(go version)"
echo ""

# 检测 Go bin 路径
GO_BIN=$(go env GOPATH)/bin
echo "📦 Go bin 路径: $GO_BIN"
echo ""

# 安装 go-start
echo "⬇️  正在安装 go-start..."
go install github.com/Martindeeepdark/go-start/cmd/go-start@latest

if [ ! -f "$GO_BIN/go-start" ]; then
    echo "❌ 安装失败"
    exit 1
fi

echo "✅ go-start 已安装到: $GO_BIN/go-start"
echo ""

# 检查是否在 PATH 中
if echo $PATH | grep -q "$GO_BIN"; then
    echo "✅ $GO_BIN 已在 PATH 中"
    echo ""
    echo "🎉 安装完成!可以直接使用:"
    echo "   go-start --version"
    echo ""
else
    echo "⚠️  $GO_BIN 不在 PATH 中,正在自动添加..."
    echo ""

    # 检测 shell 类型
    SHELL_RC=""
    if [ -n "$ZSH_VERSION" ] || [ -f ~/.zshrc ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ] || [ -f ~/.bashrc ]; then
        SHELL_RC="$HOME/.bashrc"
    elif [ -f ~/.profile ]; then
        SHELL_RC="$HOME/.profile"
    fi

    if [ -n "$SHELL_RC" ]; then
        # 添加到配置文件
        echo '' >> "$SHELL_RC"
        echo '# Go bin path' >> "$SHELL_RC"
        echo "export PATH=\"\$PATH:$GO_BIN\"" >> "$SHELL_RC"

        echo "✅ 已添加到: $SHELL_RC"
        echo ""
        echo "📝 请运行以下命令使配置生效:"
        echo "   source $SHELL_RC"
        echo ""
        echo "   或者重新打开终端"
        echo ""
        echo "   然后就可以使用:"
        echo "   go-start --version"
        echo ""
    else
        echo "⚠️  无法自动检测 shell 配置文件"
        echo ""
        echo "📝 请手动运行以下命令:"
        echo "   export PATH=\"\$PATH:$GO_BIN\""
        echo ""
        echo "   或者添加到你的 shell 配置文件中"
        echo ""
    fi
fi

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                                                           ║"
echo "║   ✅ 安装完成!                                            ║"
echo "║                                                           ║"
echo "║   快速开始:                                                ║"
echo "║   go-start create my-api                                 ║"
echo "║   go-start create --wizard                               ║"
echo "║                                                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
