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
)

//go:embed templates
var templatesFS embed.FS

var module string

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <项目名称>",
		Short: "创建新的 MVC 项目",
		Long: `创建一个基于 Gin 框架的 Go Web 项目（MVC 架构）。

示例:
  go-start create my-api
  go-start create my-api --module=github.com/yourname/my-api`,
		Args: cobra.ExactArgs(1),
		RunE: runCreate,
	}

	cmd.Flags().StringVarP(&module, "module", "m", "", "Go 模块名 (默认: 自动检测或使用项目名)")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	if !isValidProjectName(projectName) {
		return fmt.Errorf("无效的项目名称: %s", projectName)
	}

	if module == "" {
		module = detectModulePath(projectName)
		fmt.Printf("📦 模块路径: %s\n", module)
	}

	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		cwd, _ := os.Getwd()
		return fmt.Errorf("创建项目目录失败: %w\n当前目录: %s", err, cwd)
	}

	if err := generateMVCProject(projectDir, projectName, module); err != nil {
		return err
	}

	fmt.Printf("\n\033[32m✓ 项目 %s 创建成功！\033[0m\n\n", projectName)
	fmt.Println("下一步:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  cp config.yaml.example config.yaml")
	fmt.Println("  # 编辑 config.yaml 配置数据库")
	fmt.Println("  go run cmd/server/main.go")

	return nil
}

type templateData struct {
	ProjectName string
	Module      string
	Database    string
	WithRedis   bool
	ServerPort  int
}

func generateMVCProject(projectDir, projectName, mod string) error {
	dirs := []string{
		"cmd/server",
		"internal/controller",
		"internal/service",
		"internal/repository",
		"internal/model",
		"config",
		"pkg/database",
		"pkg/httpx/response",
		"pkg/httpx/router",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}

	data := templateData{
		ProjectName: projectName,
		Module:      mod,
		Database:    "mysql",
		WithRedis:   false,
		ServerPort:  8080,
	}

	if err := generateGoMod(projectDir, mod); err != nil {
		return err
	}

	templateFiles := map[string]string{
		"cmd/server/main.go":                "mvc/main.go.tpl",
		"config/config.go":                  "mvc/config/config.go.tpl",
		"config.yaml.example":               "mvc/config.yaml.tpl",
		"internal/model/user.go":            "mvc/model/user.go.tpl",
		"internal/repository/user.go":       "mvc/repository/user.go.tpl",
		"internal/repository/repository.go": "mvc/repository/repository.go.tpl",
		"internal/service/user.go":          "mvc/service/user.go.tpl",
		"internal/service/service.go":       "mvc/service/service.go.tpl",
		"internal/controller/user.go":       "mvc/controller/user.go.tpl",
		"internal/controller/controller.go": "mvc/controller/controller.go.tpl",
		"pkg/database/database.go":          "mvc/pkg/database/database.go.tpl",
		"pkg/database/driver.go":            "mvc/pkg/database/driver.go.tpl",
		"pkg/httpx/response/response.go":    "mvc/pkg/httpx/response/response.go.tpl",
		"README.md":                         "mvc/README.md.tpl",
		".gitignore":                        "mvc/gitignore.tpl",
	}

	for outputPath, templateName := range templateFiles {
		if err := renderTemplate(projectDir, outputPath, templateName, data); err != nil {
			return fmt.Errorf("生成 %s 失败: %w", outputPath, err)
		}
	}

	return nil
}

func generateGoMod(projectDir, mod string) error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.27.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)
`, mod)

	return os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(content), 0644)
}

func renderTemplate(projectDir, outputPath, templateName string, data any) error {
	templateContent, err := templatesFS.ReadFile("templates/" + templateName)
	if err != nil {
		// fallback: try reading from disk (development mode)
		fallbackPath := filepath.Join("templates", templateName)
		templateContent, err = os.ReadFile(fallbackPath)
		if err != nil {
			return fmt.Errorf("模板 %s 不存在: %w", templateName, err)
		}
	}

	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	fullPath := filepath.Join(projectDir, outputPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}

	fmt.Printf("  ✓ %s\n", outputPath)
	return nil
}

func copyEmbedDir(src string, dst string) error {
	entries, err := fs.ReadDir(templatesFS, src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := src + "/" + entry.Name()
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyEmbedDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := templatesFS.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func isValidProjectName(name string) bool {
	return name != "" && !strings.ContainsAny(name, "/\\")
}

func detectModulePath(projectName string) string {
	if parent := getParentModulePath(); parent != "" {
		return parent + "/" + projectName
	}
	return "github.com/yourname/" + projectName
}

func getParentModulePath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			if modulePath := extractModulePath(goModPath); modulePath != "" {
				return modulePath
			}
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}
	return ""
}

func extractModulePath(goModPath string) string {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "module ")), `"`)
		}
	}
	return ""
}
