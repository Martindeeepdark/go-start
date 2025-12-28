package main

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"

    "github.com/spf13/cobra"
)

// newDoctorCmd 创建 doctor 命令
// 用于检查本地开发环境与项目配置的常见问题，并提供修复建议。
func newDoctorCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "doctor",
        Short: "诊断本地环境与项目配置",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("🔍 环境与项目诊断")

            // 1. Go 版本
            fmt.Printf("• Go 版本: %s\n", runtime.Version())

            // 2. 工作区 (go.work) 检查
            if hasGoWork() {
                fmt.Println("• go.work: ✅ 已检测到工作区配置")
            } else {
                fmt.Println("• go.work: ⚠️ 未检测到，建议在 common 与 go-start 的父目录使用 go work 管理本地联动")
                fmt.Println("  参考: go work init && go work use ./go-start ./common")
            }

            // 3. go.mod 模块路径一致性
            modPath, err := readModulePath()
            if err != nil {
                fmt.Printf("• go.mod: ❌ 读取失败: %v\n", err)
            } else {
                fmt.Printf("• go.mod: 模块路径为 %s\n", modPath)
            }

            // 4. 常见依赖提示
            fmt.Println("• 依赖建议: 建议引入 golangci-lint 与 CI 测试覆盖率，提升代码质量")
            fmt.Println("• 适配建议: 使用构建标签启用 common 集成 (-tags common_integration)，便于能力按需加载")

            fmt.Println("\n✅ 诊断完成")
            return nil
        },
    }
    return cmd
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