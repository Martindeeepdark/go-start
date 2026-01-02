package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Martindeeepdark/go-start/pkg/spec"
	"github.com/spf13/cobra"
)

var (
	specFile  string
	specDir   string
	outputDir string
)

func newSpecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spec",
		Short: "基于规范文件生成代码",
		Long: `使用 YAML 规范文件自动生成 Go 代码。

支持的功能：
  - 从 YAML 规范生成数据模型
  - 自动生成 Repository、Service、Controller
  - 生成路由注册代码
  - 生成请求验证器

示例：
  # 从单个规范文件生成代码
  go-start spec generate --file=blog.spec.yaml

  # 从目录批量生成
  go-start spec generate --dir=./specs

  # 验证规范文件
  go-start spec validate --file=blog.spec.yaml

  # 创建规范文件示例
  go-start spec init`,
	}

	cmd.AddCommand(newSpecGenerateCmd())
	cmd.AddCommand(newSpecValidateCmd())
	cmd.AddCommand(newSpecInitCmd())

	return cmd
}

func newSpecGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "从规范文件生成代码",
		Long:  "从 YAML 规范文件自动生成 Go 代码",
		RunE:  runSpecGenerate,
	}

	cmd.Flags().StringVarP(&specFile, "file", "f", "", "规范文件路径")
	cmd.Flags().StringVarP(&specDir, "dir", "d", "", "规范文件目录（批量生成）")
	cmd.Flags().StringVarP(&outputDir, "output", "o", ".", "输出目录")

	return cmd
}

func newSpecValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "验证规范文件",
		Long:  "验证 YAML 规范文件的格式和内容",
		RunE:  runSpecValidate,
	}

	cmd.Flags().StringVarP(&specFile, "file", "f", "", "规范文件路径（必填）")

	return cmd
}

func newSpecInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "创建规范文件示例",
		Long:  "在当前目录创建一个规范文件示例",
		RunE:  runSpecInit,
	}

	return cmd
}

func runSpecGenerate(cmd *cobra.Command, args []string) error {
	// 检查参数
	if specFile == "" && specDir == "" {
		return fmt.Errorf("请指定 --file 或 --dir 参数")
	}

	if specFile != "" && specDir != "" {
		return fmt.Errorf("--file 和 --dir 不能同时使用")
	}

	parser := spec.New("")

	// 解析规范文件
	if specFile != "" {
		// 单个文件
		fmt.Printf("📄 正在解析规范文件: %s\n", specFile)

		s, err := parser.ParseFile(specFile)
		if err != nil {
			return fmt.Errorf("解析规范文件失败: %w", err)
		}

		// 生成代码
		generator := spec.NewGenerator(s, outputDir)
		if err := generator.Generate(); err != nil {
			return err
		}

		fmt.Printf("\n📊 生成统计:\n")
		fmt.Printf("  模型数量: %d\n", len(s.Models))
		fmt.Printf("  API 数量: %d\n", len(s.APIs))
		fmt.Printf("  验证器数量: %d\n", len(s.Requests))

	} else {
		// 批量处理目录
		fmt.Printf("📁 正在解析目录: %s\n", specDir)

		specs, err := parser.ParseDir(specDir)
		if err != nil {
			return fmt.Errorf("解析目录失败: %w", err)
		}

		if len(specs) == 0 {
			return fmt.Errorf("目录中没有找到规范文件")
		}

		fmt.Printf("找到 %d 个规范文件\n\n", len(specs))

		// 生成所有规范
		for i, s := range specs {
			fmt.Printf("[%d/%d] 生成 %s...\n", i+1, len(specs), s.Name)

			generator := spec.NewGenerator(s, outputDir)
			if err := generator.Generate(); err != nil {
				return fmt.Errorf("生成 %s 失败: %w", s.Name, err)
			}
			fmt.Println()
		}
	}

	fmt.Printf("\n✅ 代码生成完成！\n")
	fmt.Printf("📂 输出目录: %s\n\n", outputDir)

	fmt.Println("🚀 下一步:")
	fmt.Println("  1. 查看生成的代码")
	fmt.Println("  2. 根据需要调整业务逻辑")
	fmt.Println("  3. 在 main.go 中注册路由")
	fmt.Println("  4. 运行 go mod tidy")
	fmt.Println("  5. 启动服务测试")

	return nil
}

func runSpecValidate(cmd *cobra.Command, args []string) error {
	if specFile == "" {
		return fmt.Errorf("请使用 --file 参数指定规范文件")
	}

	fmt.Printf("🔍 正在验证规范文件: %s\n\n", specFile)

	parser := spec.New("")
	s, err := parser.ParseFile(specFile)
	if err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
		return err
	}

	fmt.Println("✅ 规范文件验证通过！")
	fmt.Println("📊 规范信息:")
	fmt.Printf("  名称: %s\n", s.Name)
	fmt.Printf("  版本: %s\n", s.Version)
	fmt.Printf("  模块: %s\n", s.Project.Module)
	fmt.Printf("  作者: %s\n", s.Project.Author)
	fmt.Printf("\n📦 统计:")
	fmt.Printf("  模型数量: %d\n", len(s.Models))
	fmt.Printf("  端点数量: %d\n", len(s.APIs))
	fmt.Printf("  验证器数量: %d\n", len(s.Requests))
	fmt.Printf("  业务规则数量: %d\n", len(s.Rules))

	return nil
}

func runSpecInit(cmd *cobra.Command, args []string) error {
	fmt.Println("📝 创建规范文件示例...")

	// 复制示例规范文件到当前目录
	exampleSpecPath := filepath.Join(getSpecExampleDir(), "example.blog.spec.yaml")
	outputPath := "example.spec.yaml"

	// 检查文件是否存在
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("⚠️  文件 %s 已存在\n", outputPath)
		fmt.Print("是否覆盖？(y/n): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("❌ 操作已取消")
			return nil
		}
	}

	// 读取示例文件
	content, err := os.ReadFile(exampleSpecPath)
	if err != nil {
		return fmt.Errorf("读取示例文件失败: %w", err)
	}

	// 写入到当前目录
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return fmt.Errorf("创建规范文件失败: %w", err)
	}

	fmt.Printf("\n✅ 规范文件示例已创建: %s\n\n", outputPath)
	fmt.Println("📖 使用说明:")
	fmt.Println("  1. 编辑 example.spec.yaml 文件，定义你的 API")
	fmt.Println("  2. 运行 go-start spec generate --file=example.spec.yaml")
	fmt.Println("  3. 查看生成的代码")

	return nil
}

func getSpecExampleDir() string {
	// 获取规范示例文件目录
	if _, err := os.Stat("spec"); err == nil {
		// 运行从源码
		dir, _ := filepath.Abs("spec")
		return dir
	}
	// 运行从二进制
	dir, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "..", "spec"))
	return dir
}
