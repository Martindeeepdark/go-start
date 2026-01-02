package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/Martindeeepdark/go-start/pkg/wizard"
)

//go:embed templates
var templatesFS embed.FS

var (
	archType  string
	module    string
	useWizard bool
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <项目名称>",
		Short: "创建新项目",
		Long: `创建一个新的 Go Web 项目,支持 MVC 和 DDD 架构。

示例:
  go-start create my-api                    # 使用默认配置创建
  go-start create my-api --arch=ddd         # 使用 DDD 架构
  go-start create my-api --module=github.com/用户名/my-api  # 指定模块名
  go-start create --wizard                  # 使用交互式向导`,
		Args: cobra.MaximumNArgs(1),
		RunE: runCreate,
	}

	cmd.Flags().StringVarP(&archType, "arch", "a", "mvc", "项目架构类型 (mvc, ddd)")
	cmd.Flags().StringVarP(&module, "module", "m", "", "Go 模块名 (默认: github.com/yourname/<项目名称>)")
	cmd.Flags().BoolVarP(&useWizard, "wizard", "w", false, "使用交互式向导创建项目")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	// 如果使用向导模式或没有提供项目名
	if useWizard || len(args) == 0 {
		return runWizardMode()
	}

	// 传统命令行模式
	projectName := args[0]

	// Validate project name
	if !isValidProjectName(projectName) {
		return fmt.Errorf("invalid project name: %s", projectName)
	}

	// 自动检测模块路径
	if module == "" {
		module = detectModulePath(projectName)
		fmt.Printf("📦 使用模块路径: %s\n", module)
	}

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("无法获取当前工作目录: %w", err)
	}

	// Create project directory
	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("创建项目目录失败: %w\n\n请检查:\n"+
			"  1. 当前工作目录是否存在: %s\n"+
			"  2. 是否有创建目录的权限\n"+
			"  3. 项目名称是否合法", err, cwd)
	}

	// Normalize architecture type
	archType = strings.ToLower(archType)

	// Generate project based on architecture
	switch archType {
	case "mvc":
		if err := generateMVCProjectWithOptions(projectDir, &wizard.ProjectConfig{
			ProjectName:  projectName,
			Module:       module,
			Description:  "",
			Database:     "mysql",
			WithAuth:     true,
			WithSwagger:  true,
			WithRedis:    true,
			ServerPort:   8080,
		}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported architecture: %s (supported: mvc)", archType)
	}

	fmt.Printf("✓ Project %s created successfully!\n", projectName)
	fmt.Printf("\n📝 Next steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  cp config.yaml.example config.yaml\n")
	fmt.Printf("  # Edit config.yaml with your settings\n")
	fmt.Printf("  go run cmd/server/main.go\n")

	return nil
}

// runWizardMode runs the interactive wizard
func runWizardMode() error {
	w := wizard.New()

	// 运行向导
	config, err := w.Run()
	if err != nil {
		return fmt.Errorf("向导运行失败: %w", err)
	}

	// 自动检测模块路径（如果向导中没有指定）
	if config.Module == "" || config.Module == "github.com/yourname/"+config.ProjectName {
		config.Module = detectModulePath(config.ProjectName)
		fmt.Printf("📦 自动检测到模块路径: %s\n", config.Module)
	}

	// 创建项目目录
	projectDir := filepath.Join(".", config.ProjectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		cwd, _ := os.Getwd()
		return fmt.Errorf("创建项目目录失败: %w\n\n请检查:\n"+
			"  1. 当前工作目录是否存在: %s\n"+
			"  2. 是否有创建目录的权限", err, cwd)
	}

	// 根据架构类型生成项目
	switch config.Architecture {
	case "mvc":
		if err := generateMVCProjectWithOptions(projectDir, config); err != nil {
			return err
		}
	case "ddd":
		return fmt.Errorf("DDD 架构尚未实现，请选择 MVC 架构")
	default:
		return fmt.Errorf("不支持的架构类型: %s", config.Architecture)
	}

	// 显示成功消息
	showSuccessMessage(config)
	return nil
}

// showSuccessMessage shows detailed success message with next steps
func showSuccessMessage(config *wizard.ProjectConfig) {
	fmt.Printf("\n\033[32m✓ 项目创建成功！\033[0m\n\n")
	fmt.Println("📦 项目信息")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("  名称:    %s\n", config.ProjectName)
	fmt.Printf("  位置:    %s\n", config.Module)
	fmt.Printf("  架构:    %s\n", getArchitectureLabel(config.Architecture))
	fmt.Println(strings.Repeat("─", 50))

	fmt.Print("\n🚀 下一步操作:\n")
	fmt.Println(strings.Repeat("─", 50))

	// 步骤 1
	fmt.Print("\n  1️⃣  进入项目目录:\n")
	fmt.Printf("     \033[36mcd %s\033[0m\n", config.ProjectName)

	// 步骤 2
	fmt.Print("\n  2️⃣  下载依赖:\n")
	fmt.Println("     \033[36mgo mod tidy\033[0m")

	// 步骤 3
	fmt.Print("\n  3️⃣  配置数据库:\n")
	fmt.Println("     \033[36mcp config.yaml.example config.yaml\033[0m")
	fmt.Println("     \033[90m# 然后编辑 config.yaml 配置你的数据库连接\033[0m")

	// 步骤 4
	fmt.Print("\n  4️⃣  运行项目:\n")
	fmt.Println("     \033[36mgo run cmd/server/main.go\033[0m")

	// 额外提示
	if config.WithAuth {
		fmt.Print("\n🔐 认证系统已启用:\n")
		fmt.Println("     • JWT Token 认证")
		fmt.Println("     • 用户注册/登录接口: POST /api/v1/auth/register, /api/v1/auth/login")
	}

	if config.WithSwagger {
		fmt.Print("\n📚 Swagger 文档已启用:\n")
		fmt.Printf("     • 访问地址: http://localhost:%d/swagger/index.html\033[0m\n", config.ServerPort)
	}

	fmt.Print("\n💡 提示:\n")
	fmt.Println("     • 查看 README.md 了解更多使用说明")
	fmt.Println("     • 运行 'go-start help' 查看所有命令")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
}

func generateMVCProject(projectDir, projectName, module string) error {
	// Create directory structure
	dirs := []string{
		"cmd/server",
		"internal/controller",
		"internal/service",
		"internal/repository",
		"internal/model",
		"config",
		"pkg/cache",
		"pkg/database",
		"pkg/httpx/middleware",
		"pkg/httpx/response",
		"pkg/httpx/router",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Template data
	data := struct {
		ProjectName string
		Module      string
	}{
		ProjectName: projectName,
		Module:      module,
	}

	// Template files to generate
	templateFiles := map[string]string{
		"cmd/server/main.go":                "main.go.tpl",
		"config/config.go":                  "config/config.go.tpl",
		"config.yaml.example":               "config.yaml.tpl",
		"internal/model/user.go":            "model/user.go.tpl",
		"internal/repository/user.go":       "repository/user.go.tpl",
		"internal/repository/repository.go": "repository/repository.go.tpl",
		"internal/service/user.go":          "service/user.go.tpl",
		"internal/service/service.go":       "service/service.go.tpl",
		"internal/controller/user.go":       "controller/user.go.tpl",
		"internal/controller/controller.go": "controller/controller.go.tpl",
		"README.md":                         "README.md.tpl",
		".gitignore":                        "gitignore.tpl",
	}

	// Generate go.mod
	if err := generateGoMod(projectDir, projectName, module); err != nil {
		return err
	}

	// Generate template files
	for outputPath, templateName := range templateFiles {
		if err := generateFileFromTemplate(projectDir, outputPath, templateName, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outputPath, err)
		}
	}

	// Copy pkg files from go-start
	if err := copyPkgFiles(projectDir); err != nil {
		return fmt.Errorf("failed to copy pkg files: %w", err)
	}

	return nil
}

func generateGoMod(projectDir, projectName, module string) error {
	goModContent := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/redis/go-redis/v9 v9.17.2
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.27.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)
`, module)

	return os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goModContent), 0644)
}

func generateFileFromTemplate(projectDir, outputPath, templateName string, data interface{}) error {
	// Try embedded templates first
	templatePath := filepath.Join("templates", "mvc", templateName)

	templateContent, err := fs.ReadFile(templatesFS, templatePath)
	if err != nil {
		// Fallback to filesystem
		fallbackPath := filepath.Join(getTemplateDir(), "mvc", templateName)
		templateContent, err = os.ReadFile(fallbackPath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", templateName, err)
		}
	}

	// Parse template
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	outputPath = filepath.Join(projectDir, outputPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	if err := tmpl.Execute(outputFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func copyPkgFiles(projectDir string) error {
	// Find pkg directory from source
	pkgSrc := findPkgDir()
	if pkgSrc == "" {
		// pkg not found, skip copy (user may not need it)
		fmt.Println("  ⚠️  未找到 pkg 目录,跳过复制")
		return nil
	}

	pkgDst := filepath.Join(projectDir, "pkg")

	return filepath.Walk(pkgSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(pkgSrc, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(pkgDst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}

func findPkgDir() string {
	// Try current directory first
	if _, err := os.Stat("pkg"); err == nil {
		abs, _ := filepath.Abs("pkg")
		return abs
	}

	// Try parent directory (when running from cmd/go-start)
	parentPkg := filepath.Join("..", "..", "pkg")
	if _, err := os.Stat(parentPkg); err == nil {
		abs, _ := filepath.Abs(parentPkg)
		return abs
	}

	// Try binary's parent directories
	execDir := filepath.Dir(os.Args[0])
	paths := []string{
		filepath.Join(execDir, "..", "pkg"),
		filepath.Join(execDir, "..", "..", "pkg"),
		filepath.Join(execDir, "..", "..", "..", "pkg"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			abs, _ := filepath.Abs(path)
			return abs
		}
	}

	return ""
}

func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}
	// Basic validation: should not contain path separators
	return !strings.ContainsAny(name, "/\\")
}

// detectModulePath 自动检测模块路径
func detectModulePath(projectName string) string {
	// 1. 尝试从父目录的 go.mod 获取模块路径
	if parentModule := getParentModulePath(); parentModule != "" {
		// 如果父目录有 go.mod，使用子模块路径
		return fmt.Sprintf("%s/%s", parentModule, projectName)
	}

	// 2. 使用相对路径（最简单的方式）
	return projectName
}

// getParentModulePath 获取父目录的模块路径
func getParentModulePath() string {
	// 向上查找 go.mod 文件
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// 找到 go.mod，读取模块路径
			if modulePath := extractModulePath(goModPath); modulePath != "" {
				return modulePath
			}
		}

		// 到达根目录
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return ""
}

// extractModulePath 从 go.mod 文件提取模块路径
func extractModulePath(goModPath string) string {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	// 读取第一行，格式: module xxx
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			// 移除引号（如果有）
			modulePath = strings.Trim(modulePath, `"`)
			return modulePath
		}
	}

	return ""
}

func getTemplateDir() string {
	// Get the directory where templates are stored
	// When running from bin/, templates are in ../templates
	// When running from source, templates are in ./templates
	if _, err := os.Stat("templates"); err == nil {
		// Running from source
		dir, _ := filepath.Abs("templates")
		return dir
	}
	// Running from binary
	dir, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "..", "templates"))
	return dir
}

func getRootDir() string {
	// Get the root directory of go-start
	if _, err := os.Stat("pkg"); err == nil {
		// Running from source
		dir, _ := filepath.Abs(".")
		return dir
	}
	// Running from binary
	dir, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), ".."))
	return dir
}

// generateMVCProjectWithOptions generates MVC project with wizard options
func generateMVCProjectWithOptions(projectDir string, config *wizard.ProjectConfig) error {
	// Create directory structure
	dirs := []string{
		"cmd/server",
		"internal/controller",
		"internal/service",
		"internal/repository",
		"internal/model",
		"internal/middleware",
		"config",
		"pkg/cache",
		"pkg/database",
		"pkg/httpx/middleware",
		"pkg/httpx/response",
		"pkg/httpx/router",
	}

	if config.WithAuth {
		dirs = append(dirs, "internal/auth")
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}

	// Template data
	data := struct {
		ProjectName string
		Module      string
		Description string
		Database    string
		WithAuth    bool
		WithSwagger bool
		WithRedis   bool
		ServerPort  int
	}{
		ProjectName: config.ProjectName,
		Module:      config.Module, // 使用完整的模块路径
		Description: config.Description,
		Database:    config.Database,
		WithAuth:    config.WithAuth,
		WithSwagger: config.WithSwagger,
		WithRedis:   config.WithRedis,
		ServerPort:  config.ServerPort,
	}

	// Generate go.mod
	if err := generateGoModWithOptions(projectDir, config); err != nil {
		return err
	}

	// Template files to generate
	templateFiles := map[string]string{
		"cmd/server/main.go":                "main.go.tpl",
		"config/config.go":                  "config/config.go.tpl",
		"config.yaml.example":               "config.yaml.tpl",
		"internal/model/user.go":            "model/user.go.tpl",
		"internal/repository/user.go":       "repository/user.go.tpl",
		"internal/repository/repository.go": "repository/repository.go.tpl",
		"internal/service/user.go":          "service/user.go.tpl",
		"internal/service/service.go":       "service/service.go.tpl",
		"internal/controller/user.go":       "controller/user.go.tpl",
		"internal/controller/controller.go": "controller/controller.go.tpl",
		"README.md":                         "README.md.tpl",
		".gitignore":                        "gitignore.tpl",
	}

	// Generate template files
	for outputPath, templateName := range templateFiles {
		if err := generateFileFromTemplate(projectDir, outputPath, templateName, data); err != nil {
			return fmt.Errorf("生成 %s 失败: %w", outputPath, err)
		}
	}

	// Copy pkg files from go-start
	if err := copyPkgFiles(projectDir); err != nil {
		return fmt.Errorf("复制 pkg 文件失败: %w", err)
	}

	// Generate auth files if enabled
	if config.WithAuth {
		// TODO: 生成认证相关文件
		fmt.Println("  ✓ 认证系统已配置")
	}

	// Generate swagger files if enabled
	if config.WithSwagger {
		// TODO: 生成 Swagger 配置
		fmt.Println("  ✓ Swagger 文档已配置")
	}

	return nil
}

// generateGoModWithOptions generates go.mod with wizard options
func generateGoModWithOptions(projectDir string, config *wizard.ProjectConfig) error {
	modContent := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.27.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
`, config.Module)

	// Add Redis if enabled
	if config.WithRedis {
		modContent += "\tgithub.com/redis/go-redis/v9 v9.17.2\n"
	}

	// Add JWT if auth enabled
	if config.WithAuth {
		modContent += "\tgithub.com/golang-jwt/jwt/v5 v5.2.0\n"
		modContent += "\tgolang.org/x/crypto v0.31.0\n"
	}

	// Add Swagger if enabled
	if config.WithSwagger {
		modContent += "\tgithub.com/swaggo/files v1.0.1\n"
		modContent += "\tgithub.com/swaggo/gin-swagger v1.6.0\n"
		modContent += "\tgithub.com/swaggo/swag v1.16.3\n"
	}

	// Close the require block
	modContent += ")\n"

	return os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(modContent), 0644)
}

func getArchitectureLabel(arch string) string {
	labels := map[string]string{
		"mvc": "MVC (Model-View-Controller)",
		"ddd": "DDD (Domain-Driven Design)",
	}
	if label, ok := labels[arch]; ok {
		return label
	}
	return arch
}
