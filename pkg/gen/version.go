package gen

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// checkGoVersion 检查 Go 版本是否兼容
func checkGoVersion() error {
	fmt.Println("🔍 检查 Go 版本...")

	// 获取 Go 版本
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("无法获取 Go 版本: %w", err)
	}

	// 解析版本号: "go version go1.21.0 darwin/arm64"
	// 先去掉 "go version " 前缀
	outputStr := strings.TrimPrefix(string(output), "go version ")
	// 提取版本号部分: "go1.21.0"
	versionPart := strings.Fields(outputStr)[0] // 按空格分割，取第一部分
	// 去掉 "go" 前缀: "1.21.0"
	versionStr := strings.TrimPrefix(versionPart, "go")

	parts := strings.Split(versionStr, ".")
	if len(parts) < 2 {
		return fmt.Errorf("无法解析 Go 版本")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("无法解析 Go 主版本号: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("无法解析 Go 次版本号: %w", err)
	}

	// 检查版本: Go 1.21-1.23 推荐，1.24+ 可能有兼容性问题
	if major == 1 && minor >= 21 && minor <= 23 {
		fmt.Printf("     ✓ Go 版本: %d.%d (推荐)\n", major, minor)
		return nil
	}

	if major == 1 && minor >= 24 {
		fmt.Printf("     ⚠️  Go 版本: %d.%d (可能存在兼容性问题)\n", major, minor)
		fmt.Println("     💡 推荐使用 Go 1.21-1.23")
		fmt.Println("     🔗 https://github.com/golang/go/issues/69958")
		return nil
	}

	if major < 1 || (major == 1 && minor < 21) {
		fmt.Printf("     ⚠️  Go 版本: %d.%d (过低，推荐 1.21+)\n", major, minor)
		return fmt.Errorf("Go 版本过低，请升级到 1.21 或更高版本")
	}

	fmt.Printf("     ✓ Go 版本: %d.%d\n", major, minor)
	return nil
}
