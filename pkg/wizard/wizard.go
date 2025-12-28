package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Question represents a wizard question
type Question struct {
	Text        string             // 问题文本
	Options     []string           // 选项（如果有）
	Default     string             // 默认值
	Required    bool               // 是否必填
	Validator   func(string) error // 验证函数
	Placeholder string             // 占位符提示
	Hint        string             // 提示信息
}

// ProjectConfig holds the wizard configuration
type ProjectConfig struct {
	ProjectName  string
	Module       string
	Architecture string
	Database     string
	WithAuth     bool
	WithSwagger  bool
	WithRedis    bool
	ServerPort   int
	Description  string
}

// Wizard represents the interactive wizard
type Wizard struct {
	reader *bufio.Reader
}

// New creates a new wizard instance
func New() *Wizard {
	return &Wizard{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run starts the interactive wizard
func (w *Wizard) Run() (*ProjectConfig, error) {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   🚀 欢迎使用 go-start 交互式项目创建向导                  ║
║                                                           ║
║   我将帮你创建一个专业的 Go Web 项目                       ║
║   请回答以下问题来配置你的项目                             ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`)

	config := &ProjectConfig{}

	// 1. 项目名称
	if err := w.askProjectName(config); err != nil {
		return nil, err
	}

	// 2. 模块名称
	if err := w.askModuleName(config); err != nil {
		return nil, err
	}

	// 3. 项目描述
	if err := w.askProjectDescription(config); err != nil {
		return nil, err
	}

	// 4. 架构模式
	if err := w.askArchitecture(config); err != nil {
		return nil, err
	}

	// 5. 数据库类型
	if err := w.askDatabase(config); err != nil {
		return nil, err
	}

	// 6. 是否需要 Redis
	if err := w.askRedis(config); err != nil {
		return nil, err
	}

	// 7. 是否需要认证系统
	if err := w.askAuth(config); err != nil {
		return nil, err
	}

	// 8. 是否需要 Swagger 文档
	if err := w.askSwagger(config); err != nil {
		return nil, err
	}

	// 9. 服务器端口
	if err := w.askServerPort(config); err != nil {
		return nil, err
	}

	// 显示配置摘要
	w.showSummary(config)

	// 确认创建
	if err := w.confirmCreation(config); err != nil {
		return nil, err
	}

	return config, nil
}

// askProjectName asks for the project name
func (w *Wizard) askProjectName(config *ProjectConfig) error {
	fmt.Println("\n📦 步骤 1/9: 项目名称")
	fmt.Println("═════════════════════════════════════════")

	for {
		answer, err := w.ask(Question{
			Text:      "请输入项目名称（例如: my-api）",
			Required:  true,
			Validator: validateProjectName,
			Hint:      "项目名称只能包含字母、数字和连字符，且不能以连字符开头或结尾",
		})
		if err != nil {
			return err
		}

		config.ProjectName = answer
		break
	}

	return nil
}

// askModuleName asks for the Go module name
func (w *Wizard) askModuleName(config *ProjectConfig) error {
	fmt.Println("\n📦 步骤 2/9: Go 模块名称")
	fmt.Println("═════════════════════════════════════════")

	defaultModule := fmt.Sprintf("github.com/yourname/%s", config.ProjectName)

	answer, err := w.ask(Question{
		Text:     "请输入 Go 模块名称",
		Default:  defaultModule,
		Required: true,
		Hint:     "这是你的 Go 模块路径，通常格式为: github.com/用户名/项目名",
	})
	if err != nil {
		return err
	}

	config.Module = answer
	return nil
}

// askProjectDescription asks for project description
func (w *Wizard) askProjectDescription(config *ProjectConfig) error {
	fmt.Println("\n📝 步骤 3/9: 项目描述")
	fmt.Println("═════════════════════════════════════════")

	answer, err := w.ask(Question{
		Text:        "请输入项目描述（可选）",
		Required:    false,
		Placeholder: "例如: 一个简单的 RESTful API 服务",
		Hint:        "这个描述将出现在 README.md 中",
	})
	if err != nil {
		return err
	}

	config.Description = answer
	return nil
}

// askArchitecture asks for the architecture pattern
func (w *Wizard) askArchitecture(config *ProjectConfig) error {
	fmt.Println("\n🏗️  步骤 4/9: 架构模式")
	fmt.Println("═════════════════════════════════════════")
	fmt.Println("选择你的项目架构模式：")
	fmt.Println("  1️⃣  MVC (Model-View-Controller)")
	fmt.Println("     - 适合中小型项目")
	fmt.Println("     - 简单直观，易于上手")
	fmt.Println("     - 推荐：新人首选")
	fmt.Println()
	fmt.Println("  2️⃣  DDD (Domain-Driven Design)")
	fmt.Println("     - 适合大型复杂项目")
	fmt.Println("     - 领域驱动设计，业务逻辑清晰")
	fmt.Println("     - 推荐：高级工程师")

	answer, err := w.ask(Question{
		Text:     "请选择架构模式 (1 或 2)",
		Options:  []string{"1", "2", "mvc", "ddd"},
		Default:  "1",
		Required: true,
	})
	if err != nil {
		return err
	}

	// 转换答案
	switch strings.ToLower(answer) {
	case "1", "mvc":
		config.Architecture = "mvc"
	case "2", "ddd":
		config.Architecture = "ddd"
	}

	return nil
}

// askDatabase asks for the database type
func (w *Wizard) askDatabase(config *ProjectConfig) error {
	fmt.Println("\n🗄️  步骤 5/9: 数据库类型")
	fmt.Println("═════════════════════════════════════════")
	fmt.Println("选择你使用的数据库：")
	fmt.Println("  1️⃣  MySQL")
	fmt.Println("     - 最流行的开源数据库")
	fmt.Println("     - 社区活跃，资源丰富")
	fmt.Println()
	fmt.Println("  2️⃣  PostgreSQL")
	fmt.Println("     - 功能强大的开源数据库")
	fmt.Println("     - 支持 JSON、GIS 等高级特性")
	fmt.Println()
	fmt.Println("  3️⃣  SQLite")
	fmt.Println("     - 轻量级嵌入式数据库")
	fmt.Println("     - 适合小型项目或原型开发")

	answer, err := w.ask(Question{
		Text:     "请选择数据库 (1/2/3)",
		Options:  []string{"1", "2", "3", "mysql", "postgresql", "sqlite"},
		Default:  "1",
		Required: true,
	})
	if err != nil {
		return err
	}

	switch strings.ToLower(answer) {
	case "1", "mysql":
		config.Database = "mysql"
	case "2", "postgresql", "postgres":
		config.Database = "postgresql"
	case "3", "sqlite", "sqlite3":
		config.Database = "sqlite"
	}

	return nil
}

// askRedis asks if Redis is needed
func (w *Wizard) askRedis(config *ProjectConfig) error {
	fmt.Println("\n⚡ 步骤 6/9: Redis 缓存")
	fmt.Println("═════════════════════════════════════════")
	fmt.Println("Redis 是一个高性能的键值存储系统，可用于：")
	fmt.Println("  • 缓存热点数据")
	fmt.Println("  • 会话存储")
	fmt.Println("  • 分布式锁")
	fmt.Println("  • 消息队列")

	answer, err := w.ask(Question{
		Text:     "是否需要 Redis 支持？(y/n)",
		Options:  []string{"y", "n", "yes", "no"},
		Default:  "y",
		Required: true,
	})
	if err != nil {
		return err
	}

	config.WithRedis = strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes"
	return nil
}

// askAuth asks if authentication is needed
func (w *Wizard) askAuth(config *ProjectConfig) error {
	fmt.Println("\n🔐 步骤 7/9: 用户认证系统")
	fmt.Println("═════════════════════════════════════════")
	fmt.Println("是否需要内置的用户认证系统？")
	fmt.Println("  包含功能：")
	fmt.Println("  • JWT Token 认证")
	fmt.Println("  • 用户注册/登录")
	fmt.Println("  • 密码加密存储")
	fmt.Println("  • 权限控制中间件")

	answer, err := w.ask(Question{
		Text:     "是否需要认证系统？(y/n)",
		Options:  []string{"y", "n", "yes", "no"},
		Default:  "y",
		Required: true,
	})
	if err != nil {
		return err
	}

	config.WithAuth = strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes"
	return nil
}

// askSwagger asks if Swagger documentation is needed
func (w *Wizard) askSwagger(config *ProjectConfig) error {
	fmt.Println("\n📚 步骤 8/9: API 文档")
	fmt.Println("═════════════════════════════════════════")
	fmt.Println("是否需要自动生成 Swagger API 文档？")
	fmt.Println("  优势：")
	fmt.Println("  • 自动生成在线 API 文档")
	fmt.Println("  • 支持 API 测试和调试")
	fmt.Println("  • 便于前后端协作")

	answer, err := w.ask(Question{
		Text:     "是否需要 Swagger 文档？(y/n)",
		Options:  []string{"y", "n", "yes", "no"},
		Default:  "y",
		Required: true,
	})
	if err != nil {
		return err
	}

	config.WithSwagger = strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes"
	return nil
}

// askServerPort asks for the server port
func (w *Wizard) askServerPort(config *ProjectConfig) error {
	fmt.Println("\n🔌 步骤 9/9: 服务器端口")
	fmt.Println("═════════════════════════════════════════")

	answer, err := w.ask(Question{
		Text:      "请输入服务器端口号",
		Default:   "8080",
		Required:  true,
		Validator: validatePort,
		Hint:      "建议使用 1024-49151 之间的端口",
	})
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(answer)
	if err != nil {
		return fmt.Errorf("无效的端口号: %w", err)
	}

	config.ServerPort = port
	return nil
}

// showSummary displays the configuration summary
func (w *Wizard) showSummary(config *ProjectConfig) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("📋 项目配置摘要")
	fmt.Println(strings.Repeat("═", 60))

	fmt.Printf("  项目名称:        %s\n", config.ProjectName)
	fmt.Printf("  模块名称:        %s\n", config.Module)
	if config.Description != "" {
		fmt.Printf("  项目描述:        %s\n", config.Description)
	}
	fmt.Printf("  架构模式:        %s\n", getArchitectureLabel(config.Architecture))
	fmt.Printf("  数据库:          %s\n", getDatabaseLabel(config.Database))
	fmt.Printf("  Redis 缓存:      %s\n", getBoolLabel(config.WithRedis))
	fmt.Printf("  认证系统:        %s\n", getBoolLabel(config.WithAuth))
	fmt.Printf("  Swagger 文档:    %s\n", getBoolLabel(config.WithSwagger))
	fmt.Printf("  服务端口:        %d\n", config.ServerPort)

	fmt.Println(strings.Repeat("═", 60))
}

// confirmCreation asks for final confirmation
func (w *Wizard) confirmCreation(config *ProjectConfig) error {
	fmt.Println("\n✨ 准备创建项目！")

	answer, err := w.ask(Question{
		Text:     "确认创建项目？(y/n)",
		Options:  []string{"y", "n", "yes", "no"},
		Default:  "y",
		Required: true,
	})
	if err != nil {
		return err
	}

	if strings.ToLower(answer) != "y" && strings.ToLower(answer) != "yes" {
		return fmt.Errorf("用户取消创建")
	}

	return nil
}

// ask asks a single question
func (w *Wizard) ask(q Question) (string, error) {
	for {
		// 构建提示文本
		prompt := q.Text
		if q.Default != "" {
			prompt += fmt.Sprintf(" (默认: %s)", q.Default)
		}
		prompt += ": "

		// 显示提示
		fmt.Print("\033[36m➜\033[0m " + prompt)

		// 读取输入
		answer, err := w.reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("读取输入失败: %w", err)
		}

		answer = strings.TrimSpace(answer)

		// 使用默认值
		if answer == "" && q.Default != "" {
			answer = q.Default
		}

		// 验证必填
		if answer == "" && q.Required {
			if q.Hint != "" {
				fmt.Printf("\033[33m⚠️  %s\033[0m\n", q.Hint)
			}
			fmt.Printf("\033[31m✗ 此项为必填，请重新输入\033[0m\n\n")
			continue
		}

		// 如果为空且非必填，直接返回
		if answer == "" {
			return answer, nil
		}

		// 验证选项
		if len(q.Options) > 0 {
			valid := false
			for _, opt := range q.Options {
				if strings.EqualFold(answer, opt) {
					valid = true
					break
				}
			}
			if !valid {
				fmt.Printf("\033[31m✗ 无效的选项，请选择: %s\033[0m\n\n", strings.Join(q.Options, "/"))
				continue
			}
		}

		// 自定义验证
		if q.Validator != nil {
			if err := q.Validator(answer); err != nil {
				fmt.Printf("\033[31m✗ %v\033[0m\n\n", err)
				if q.Hint != "" {
					fmt.Printf("\033[33m💡 提示: %s\033[0m\n\n", q.Hint)
				}
				continue
			}
		}

		// 验证通过
		fmt.Printf("\033[32m✓\033[0m\n\n")
		return answer, nil
	}
}

// Helper functions

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

func getDatabaseLabel(db string) string {
	labels := map[string]string{
		"mysql":      "MySQL",
		"postgresql": "PostgreSQL",
		"sqlite":     "SQLite",
	}
	if label, ok := labels[db]; ok {
		return label
	}
	return db
}

func getBoolLabel(b bool) string {
	if b {
		return "\033[32m✓ 启用\033[0m"
	}
	return "\033[90m✗ 禁用\033[0m"
}

// Validators

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("项目名称不能为空")
	}

	// 检查非法字符
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("项目名称不能包含路径分隔符")
	}

	// 检查是否以连字符开头或结尾
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("项目名称不能以连字符开头或结尾")
	}

	// 检查是否只包含有效字符
	for _, c := range name {
		if !isAlphaNumeric(c) && c != '-' && c != '_' {
			return fmt.Errorf("项目名称只能包含字母、数字、连字符和下划线")
		}
	}

	return nil
}

func validatePort(port string) error {
	p, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("端口号必须是数字")
	}

	if p < 1 || p > 65535 {
		return fmt.Errorf("端口号必须在 1-65535 之间")
	}

	return nil
}

func isAlphaNumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
