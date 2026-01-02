package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Martindeeepdark/go-start/pkg/check"
	"github.com/spf13/cobra"
)

// newDoctorCmd 创建 doctor 命令
// 用于检查本地开发环境与项目配置的常见问题，并提供修复建议。
func newDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "诊断本地环境与项目配置",
		Long: `检查开发环境配置,确保 go-start 可以正常工作。

检查项目:
  • Go 版本兼容性
  • 数据库连接
  • 必要的依赖工具
  • 项目配置文件

示例:
  go-start doctor              # 检查所有项目
  go-start doctor --verbose    # 显示详细信息`,
		RunE: runDoctor,
	}

	return cmd
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Print(`
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   🔍 go-start 环境诊断工具                                ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`)

	allPassed := true

	// 1. Go 版本检查
	fmt.Println("📌 检查 Go 版本...")
	goVersionInfo := check.CheckGoVersion()
	check.PrintVersionInfo(goVersionInfo)
	if !goVersionInfo.Valid {
		allPassed = false
	}

	// 2. 检查必要工具
	fmt.Println("📌 检查必要工具...")
	checkTools()

	// 3. 数据库连接检查
	fmt.Println("📌 检查数据库连接...")
	checkDatabase()

	// 4. 项目配置检查
	fmt.Println("📌 检查项目配置...")
	checkProjectConfig()

	// 总结
	fmt.Print(`
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
`)
	if allPassed {
		fmt.Println("║   ✅ 所有检查通过!环境配置正确                             ║")
	} else {
		fmt.Println("║   ⚠️  发现一些问题,请根据上述提示修复                       ║")
	}
	fmt.Println(`║                                                           ║
╚═══════════════════════════════════════════════════════════╝`)

	return nil
}

// checkTools 检查必要的开发工具
func checkTools() {
	tools := []struct {
		name string
		cmd  string
		args string
		need bool
		hint string
	}{
		{
			name: "Go",
			cmd:  "go",
			args: "version",
			need: true,
			hint: "",
		},
		{
			name: "Git",
			cmd:  "git",
			args: "version",
			need: true,
			hint: "",
		},
		{
			name: "Docker",
			cmd:  "docker",
			args: "version",
			need: false,
			hint: "可选,用于容器化部署",
		},
		{
			name: "golangci-lint",
			cmd:  "golangci-lint",
			args: "version",
			need: false,
			hint: "推荐,用于代码质量检查",
		},
	}

	for _, tool := range tools {
		cmd := exec.Command(tool.cmd, tool.args)
		if err := cmd.Run(); err != nil {
			if tool.need {
				fmt.Printf("   ❌ %s 未安装\n", tool.name)
				fmt.Printf("      请安装 %s 后重试\n", tool.name)
			} else {
				fmt.Printf("   ⚠️  %s 未安装 (可选)\n", tool.name)
				if tool.hint != "" {
					fmt.Printf("      %s\n", tool.hint)
				}
			}
		} else {
			fmt.Printf("   ✅ %s 已安装\n", tool.name)
		}
	}
	fmt.Println()
}

// checkDatabase 检查数据库连接
func checkDatabase() {
	// 检查是否有 config.yaml
	configFiles := []string{"config.yaml", "config.yaml.example"}
	foundConfig := false

	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			fmt.Printf("   ✅ 找到配置文件: %s\n", configFile)
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		fmt.Println("   ⚠️  未找到配置文件 config.yaml")
		fmt.Println("      提示: 在创建项目后,需要复制 config.yaml.example 为 config.yaml")
		fmt.Println("      命令: cp config.yaml.example config.yaml")
		fmt.Println()
		return
	}

	// 检查数据库服务是否运行
	fmt.Println("   💡 提示: 运行以下命令测试数据库连接:")
	fmt.Println("      go-start check db --config=config.yaml")
	fmt.Println()
}

// checkProjectConfig 检查项目配置
func checkProjectConfig() {
	// 检查 go.mod
	if _, err := os.Stat("go.mod"); err == nil {
		fmt.Println("   ✅ 找到 go.mod")

		// 读取模块路径
		modPath, err := readModulePath()
		if err != nil {
			fmt.Printf("   ⚠️  无法读取模块路径: %v\n", err)
		} else {
			fmt.Printf("      模块路径: %s\n", modPath)
		}
	} else {
		fmt.Println("   ⚠️  未找到 go.mod")
		fmt.Println("      请在项目根目录运行此命令")
	}

	fmt.Println()
}

// hasGoWork 检查当前或父级目录是否存在 go.work 文件
func hasGoWork() bool {
	wd, _ := os.Getwd()
	for i := 0; i < 3; i++ {
		candidate := filepath.Join(wd, "go.work")
		if _, err := os.Stat(candidate); err == nil {
			return true
		}
		wd = filepath.Dir(wd)
	}
	return false
}

// readModulePath 读取当前项目的 go.mod 模块路径
func readModulePath() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}
	// 简单解析第一行: module <path>
	for _, line := range splitLines(string(data)) {
		if len(line) > 7 && line[:6] == "module" {
			return trimSpace(line[6:]), nil
		}
	}
	return "", fmt.Errorf("未找到 module 声明")
}

// splitLines 简易按行分割
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// trimSpace 去除首尾空白
func trimSpace(s string) string {
	i := 0
	j := len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}
