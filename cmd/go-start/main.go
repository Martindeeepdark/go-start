package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set by build flags or git tag
	// If not set, defaults to "dev"
	Version = "v1.3.1"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-start",
		Short: "快速创建 Go Web 项目的脚手架工具",
		Long: `go-start 是一个命令行工具,帮助你快速创建基于 Gin 框架的 Go Web 项目
支持 MVC 和 DDD 两种架构模式,可以从数据库自动生成完整的 CRUD API。`,
		Version: Version,
	}

	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newSpecCmd())
	rootCmd.AddCommand(newGenCmd())
	rootCmd.AddCommand(newDoctorCmd())
	rootCmd.AddCommand(newVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本号",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("go-start 版本 %s\n", Version)
			fmt.Println("GitHub: https://github.com/Martindeeepdark/go-start")
			fmt.Println("文档: https://github.com/Martindeeepdark/go-start#readme")
		},
	}
}
