#!/bin/bash

# go-start 一键安装脚本
# 适用于 macOS 和 Linux

set -e

# 避免Go toolchain问题
export GOTOOLCHAIN=local

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

# 测试 Go 是否正常工作
if ! go version &> /dev/null; then
    echo "⚠️  检测到 Go 但运行异常"
    echo "   这可能是 Go toolchain 缓存问题"
    echo ""
    echo "📝 正在尝试从源码编译安装..."
    echo ""

    # 从源码编译
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    if git clone https://github.com/Martindeeepdark/go-start.git 2>/dev/null; then
        cd go-start
        if go build -o go-start ./cmd/go-start 2>/dev/null; then
            # 安装到系统路径
            if [ -w /usr/local/bin ]; then
                mv go-start /usr/local/bin/
                GO_BIN="/usr/local/bin"
            else
                mkdir -p "$HOME/go/bin"
                mv go-start "$HOME/go/bin/"
                GO_BIN="$HOME/go/bin"
            fi

            echo "✅ 从源码编译成功!"
            echo ""
            cd "$TMP_DIR"
            rm -rf go-start

            # 继续配置 PATH
            goto_check_path
        fi
    fi

    echo "❌ 从源码编译也失败了"
    echo ""
    echo "请手动解决 Go 环境问题:"
    echo "  1. 清理缓存: rm -rf ~/go/pkg/mod/golang.org/toolchain*"
    echo "  2. 重装 Go: brew reinstall go"
    exit 1
fi

echo "✅ 检测到 Go: $(go version)"
echo ""

# 检测 Go bin 路径
GOPATH=$(go env GOPATH)
GOBIN="$GOPATH/bin"

# Go 1.24+ 会安装到平台特定子目录 (如 darwin_arm64)
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
PLATFORM_DIR="$GOBIN/${GOOS}_${GOARCH}"

# 检测实际的可执行文件路径
if [ -f "$PLATFORM_DIR/go-start" ]; then
    # Go 1.24+: 文件在平台子目录中
    ACTUAL_BIN_PATH="$PLATFORM_DIR/go-start"
    # 创建符号链接到 bin 目录
    if [ ! -L "$GOBIN/go-start" ] && [ ! -f "$GOBIN/go-start" ]; then
        ln -s "$PLATFORM_DIR/go-start" "$GOBIN/go-start" 2>/dev/null || true
    fi
else
    # Go < 1.24: 文件直接在 bin 目录中
    ACTUAL_BIN_PATH="$GOBIN/go-start"
fi

echo "📦 Go bin 路径: $GOBIN"
echo ""

# 安装 go-start
echo "⬇️  正在安装 go-start..."
echo ""

if go install github.com/Martindeeepdark/go-start/cmd/go-start@latest 2>&1; then
    echo "✅ 安装命令执行成功"
else
    echo "❌ 安装命令执行失败"
    echo ""
    echo "可能的原因:"
    echo "  1. 网络问题,无法访问 GitHub"
    echo "  2. Go 版本不兼容 (需要 Go 1.21+)"
    echo "  3. Go module 缓存问题"
    echo ""
    echo "尝试手动安装:"
    echo "  git clone https://github.com/Martindeeepdark/go-start.git"
    echo "  cd go-start"
    echo "  go build -o go-start ./cmd/go-start"
    echo "  sudo mv go-start /usr/local/bin/"
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

echo "✅ go-start 已安装到: $ACTUAL_BIN_PATH"
echo ""

# 检查是否在 PATH 中
if echo $PATH | grep -q "$GOBIN"; then
    echo "✅ $GOBIN 已在 PATH 中"
    echo ""
    echo "🎉 安装完成!可以直接使用:"
    echo "   go-start --version"
    echo ""
else
    echo "⚠️  $GOBIN 不在 PATH 中,正在自动添加..."
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
        echo "   export PATH=\"\$PATH:$GOBIN\""
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
