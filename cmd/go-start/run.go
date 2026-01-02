package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "一键启动项目 (自动配置并运行)",
		Long: `一键启动项目，自动检查配置文件、依赖，然后启动服务。

功能特性:
  1. 检查配置文件是否存在 (.env 或 config.yaml)
  2. 自动执行 go mod tidy 下载依赖
  3. 支持热加载 (如果安装了 air)
  4. 智能查找 main.go 入口文件

使用前提:
  - 确保已经配置好数据库连接 (.env 或 config.yaml)
  - 确保数据库服务已启动

示例:
  go-start run                 # 一键启动项目
  go-start run --verbose       # 显示详细日志

提示:
  - 安装 air 以支持热加载: go install github.com/cosmtrek/air@latest
  - 首次运行前，请复制配置文件: cp .env.example .env`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRun(cmd, args, verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "显示详细日志")
	return cmd
}

func runRun(cmd *cobra.Command, args []string, verbose bool) error {
	// 1. 检查是否在 Go 项目目录中
	if !isGoProject() {
		return fmt.Errorf("❌ 当前目录不是 Go 项目目录 (缺少 go.mod 文件)")
	}
	fmt.Println("✅ 找到 Go 项目")

	// 2. 检查配置文件
	if !checkConfigFiles() {
		fmt.Println("⚠️  警告: 未找到配置文件 (.env 或 config.yaml)")
		fmt.Println("💡 提示: 请先复制配置文件: cp .env.example .env")
		fmt.Println("         或: cp config.yaml.example config.yaml")
		fmt.Println("⏳  继续启动，但可能需要手动配置环境变量...")
	} else {
		fmt.Println("✅ 找到配置文件")
	}

	// 3. 下载依赖
	fmt.Println("\n📦 正在检查并下载依赖...")
	if err := goModTidy(); err != nil {
		return fmt.Errorf("❌ go mod tidy 失败: %w", err)
	}
	fmt.Println("✅ 依赖下载完成")

	// 4. 启动服务
	fmt.Println("\n🚀 准备启动服务...")
	if hasCommand("air") {
		if verbose {
			fmt.Println("🔥 使用热加载模式运行 (air)...")
			fmt.Println("💡 提示: 代码修改会自动重启服务\n")
		} else {
			fmt.Println("🔥 使用热加载模式运行 (air)...\n")
		}
		return runWithAir(verbose)
	}

	if verbose {
		fmt.Println("▶️  运行项目 (无热加载)")
		fmt.Println("💡 提示: 安装 air 以支持热加载: go install github.com/cosmtrek/air@latest\n")
	} else {
		fmt.Println("▶️  运行项目 (无热加载)\n")
	}
	return runDirectly(verbose)
}

func runWithAir(verbose bool) error {
	args := []string{}
	if !verbose {
		args = append(args, "-q") // air quiet mode
	}
	cmd := exec.Command("air", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDirectly(verbose bool) error {
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
		return fmt.Errorf("❌ 未找到 main.go 文件 (尝试了: %v)", mainPaths)
	}

	args := []string{"run", mainPath}
	if verbose {
		// go run 的详细输出
	} else {
		// 可以在这里添加其他参数
	}

	cmd := exec.Command("go", args...)
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

// checkConfigFiles 检查配置文件是否存在
func checkConfigFiles() bool {
	configFiles := []string{
		".env",
		"config.yaml",
		"config.yml",
		".env.local",
		".env.development",
	}

	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			return true
		}
	}
	return false
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
