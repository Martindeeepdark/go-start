package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "在开发模式下运行项目",
		Long: `使用 air 热加载功能运行项目。

如果未安装 air,则直接运行项目(无热加载)。

示例:
  go-start run                 # 运行项目
  go-start run --verbose       # 显示详细日志

提示:
  安装 air 以支持热加载: go install github.com/cosmtrek/air@latest`,
		RunE: runRun,
	}

	return cmd
}

func runRun(cmd *cobra.Command, args []string) error {
	// Check if we're in a Go project
	if !isGoProject() {
		return fmt.Errorf("not in a Go project directory")
	}

	// Ensure dependencies are downloaded
	fmt.Println("📦 检查并下载依赖...")
	if err := goModTidy(); err != nil {
		return fmt.Errorf("go mod tidy 失败: %w", err)
	}

	// Check if air is installed for hot reload
	if hasCommand("air") {
		fmt.Println("🔥 使用热加载模式运行 (air)...")
		return runWithAir()
	}

	// Run directly
	fmt.Println("▶️  运行项目 (无热加载)")
	fmt.Println("💡 提示: 安装 air 以支持热加载: go install github.com/cosmtrek/air@latest")
	return runDirectly()
}

func runWithAir() error {
	cmd := exec.Command("air")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDirectly() error {
	// Try to find main.go
	mainPaths := []string{
		"cmd/server/main.go",
		"main.go",
	}

	var mainPath string
	for _, path := range mainPaths {
		if _, err := os.Stat(path); err == nil {
			mainPath = path
			break
		}
	}

	if mainPath == "" {
		return fmt.Errorf("main.go not found (tried: %v)", mainPaths)
	}

	cmd := exec.Command("go", "run", mainPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isGoProject() bool {
	if _, err := os.Stat("go.mod"); err != nil {
		return false
	}
	return true
}

func hasCommand(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return false
	}
	return cmd.ProcessState.Success()
}

func goModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
