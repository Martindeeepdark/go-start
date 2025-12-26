package gen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GenerateGoMod 生成 go.mod 文件
func (g *DatabaseGenerator) GenerateGoMod() error {
	modulePath := getModulePath(g.config.Module)
	if modulePath == "github.com/yourname/project" {
		// 如果使用默认值，给出警告但仍然生成
		fmt.Println("⚠️  使用默认模块路径 github.com/yourname/project")
		fmt.Println("    建议使用 --module 参数指定您的项目路径")
		fmt.Println()
	}

	outputPath := g.config.Output
	if outputPath == "./internal" {
		outputPath = "."
	}

	// 创建 go.mod 文件
	goModPath := filepath.Join(outputPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		fmt.Println("⚠️  go.mod 已存在，跳过生成")
		return nil
	}

	fmt.Println("📦 正在生成 go.mod...")

	// 写入 go.mod 内容
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.27.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)

require (
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`, modulePath)

	if err := os.WriteFile(goModPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 go.mod 失败: %w", err)
	}

	fmt.Println("     ✓ go.mod 创建成功")
	fmt.Printf("     模块路径: %s\n", modulePath)

	// 尝试运行 go mod tidy，但不阻塞
	fmt.Println("📦 正在运行 go mod tidy (可能需要几分钟)...")
	go func() {
		if err := g.runGoModTidy(outputPath); err != nil {
			fmt.Printf("⚠️  go mod tidy 失败: %v\n", err)
			fmt.Println("     提示: 请在项目目录中手动运行 'go mod tidy'")
		} else {
			fmt.Println("     ✓ go mod tidy 完成")
		}
	}()

	return nil
}

// runGoModTidy 运行 go mod tidy
func (g *DatabaseGenerator) runGoModTidy(workDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = workDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", strings.TrimSpace(string(output)), err)
	}

	return nil
}
