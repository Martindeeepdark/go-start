package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	archType string
	module   string
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <project-name>",
		Short: "Create a new project",
		Long: `Create a new Go project with MVC architecture.

Example:
  go-start create my-api
  go-start create my-api --arch=Mvc
  go-start create my-api --module=github.com/myname/my-api`,
		Args: cobra.ExactArgs(1),
		RunE: runCreate,
	}

	cmd.Flags().StringVarP(&archType, "arch", "a", "mvc", "Project architecture (mvc, ddd)")
	cmd.Flags().StringVarP(&module, "module", "m", "", "Go module name (default: github.com/yourname/<project-name>)")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Validate project name
	if !isValidProjectName(projectName) {
		return fmt.Errorf("invalid project name: %s", projectName)
	}

	// Set default module name
	if module == "" {
		module = fmt.Sprintf("github.com/yourname/%s", projectName)
	}

	// Create project directory
	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Normalize architecture type
	archType = strings.ToLower(archType)

	// Generate project based on architecture
	switch archType {
	case "mvc":
		if err := generateMVCProject(projectDir, projectName, module); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported architecture: %s (supported: mvc)", archType)
	}

	fmt.Printf("✓ Project %s created successfully!\n", projectName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  cp config.yaml.example config.yaml\n")
	fmt.Printf("  # Edit config.yaml with your settings\n")
	fmt.Printf("  go run cmd/server/main.go\n")

	return nil
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

go 1.25.4

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
	// Get template path
	templatePath := filepath.Join(getTemplateDir(), "mvc", templateName)

	// Read template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
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
	// Copy pkg directory
	pkgSrc := filepath.Join(getRootDir(), "pkg")
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

func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}
	// Basic validation: should not contain path separators
	return !strings.ContainsAny(name, "/\\")
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
