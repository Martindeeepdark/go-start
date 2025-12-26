package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourname/go-start/pkg/gen"
)

var (
	genDSN         string
	genTables      string
	genOutput      string
	genSQLFile     string
	genInteractive bool
	genConfig      string
)

func newGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "从数据库生成 CRUD 代码",
		Long: `自动生成完整的 CRUD 代码，让你专注于业务逻辑。

支持三种生成模式：
  1. 从现有数据库生成 (gen db)
  2. 从 SQL 文件生成 (gen sql)
  3. 从 spec 文件生成 (spec generate)

生成的代码包括：
  - Model (数据模型)
  - Repository (数据访问层，包含 CRUD + 高级查询)
  - Service (业务逻辑层)
  - Controller (HTTP 处理器)
  - Routes (路由注册)

示例：
  # 交互式选择表（推荐）
  go-start gen db --dsn="root:pass@tcp(localhost:3306)/mydb" --interactive

  # 指定表名生成
  go-start gen db --dsn="..." --tables=users,articles,comments

  # 使用通配符
  go-start gen db --dsn="..." --tables="user*"

  # 从 SQL 文件生成
  go-start gen sql --file=schema.sql`,
	}

	cmd.AddCommand(newGenDbCmd())
	cmd.AddCommand(newGenSqlCmd())

	return cmd
}

// gen db - 从数据库生成
func newGenDbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "从数据库表生成 CRUD 代码",
		Long:  "连接数据库，读取表结构，生成完整的 CRUD 代码",
		RunE:  runGenDb,
	}

	cmd.Flags().StringVar(&genDSN, "dsn", "", "数据库连接字符串 (必填)")
	cmd.Flags().StringVar(&genTables, "tables", "", "要生成的表名，逗号分隔 (如: users,articles)，支持通配符 (user*)")
	cmd.Flags().BoolVar(&genInteractive, "interactive", false, "交互式选择表（推荐）")
	cmd.Flags().StringVar(&genConfig, "config", "", "从配置文件读取表列表")
	cmd.Flags().StringVar(&genOutput, "output", "./internal", "输出目录")

	return cmd
}

// gen sql - 从 SQL 文件生成
func newGenSqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sql",
		Short: "从 SQL DDL 文件生成 CRUD 代码",
		Long:  "解析 SQL DDL 文件，生成完整的 CRUD 代码",
		RunE:  runGenSql,
	}

	cmd.Flags().StringVar(&genSQLFile, "file", "", "SQL 文件路径 (必填)")
	cmd.Flags().StringVar(&genOutput, "output", "./internal", "输出目录")

	return cmd
}

func runGenDb(cmd *cobra.Command, args []string) error {
	// 验证参数
	if genDSN == "" {
		return fmt.Errorf("请提供数据库连接字符串 (--dsn)")
	}

	// 解析要生成的表列表
	var tables []string
	var err error

	if genInteractive {
		// 交互式模式
		tables, err = selectTablesInteractive(genDSN)
		if err != nil {
			return err
		}
		if len(tables) == 0 {
			fmt.Println("❌ 未选择任何表，操作已取消")
			return nil
		}
	} else if genConfig != "" {
		// 从配置文件读取
		tables, err = loadTablesFromConfig(genConfig)
		if err != nil {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	} else if genTables != "" {
		// 命令行指定
		tables = parseTables(genTables)
	} else {
		// 未指定，提示用户
		return fmt.Errorf("请使用以下方式之一指定要生成的表：\n" +
			"  1. --tables=users,articles (指定表名)\n" +
			"  2. --interactive (交互式选择，推荐)\n" +
			"  3. --config=gen.yaml (从配置文件读取)")
	}

	fmt.Printf("\n🔌 正在连接数据库...\n")
	fmt.Printf("📊 DSN: %s\n", maskDSN(genDSN))
	fmt.Printf("📋 将生成 %d 张表: %s\n\n", len(tables), strings.Join(tables, ", "))

	// 创建生成器
	generator := gen.NewDatabaseGenerator(gen.Config{
		DSN:    genDSN,
		Tables: tables,
		Output: genOutput,
	})

	// 生成代码
	if err := generator.Generate(); err != nil {
		return fmt.Errorf("生成代码失败: %w", err)
	}

	fmt.Println("\n✅ 代码生成完成！")
	fmt.Println("\n📦 已生成:")
	fmt.Println("  ✓ Model (数据模型)")
	fmt.Println("  ✓ Repository (数据访问层 + CRUD + 高级查询)")
	fmt.Println("  ✓ Service (业务逻辑层 + 缓存)")
	fmt.Println("  ✓ Controller (HTTP 处理器 + RESTful API)")
	fmt.Println("  ✓ Routes (路由注册)")

	fmt.Println("\n🚀 下一步:")
	fmt.Println("  1. 检查生成的代码")
	fmt.Println("  2. 在 Service 层添加自定义业务逻辑")
	fmt.Println("  3. 在 main.go 中注册路由: import internal/routes")
	fmt.Println("  4. 运行 go mod tidy")
	fmt.Println("  5. 启动服务: go run cmd/server/main.go")

	return nil
}

func runGenSql(cmd *cobra.Command, args []string) error {
	if genSQLFile == "" {
		return fmt.Errorf("请提供 SQL 文件路径 (--file)")
	}

	if _, err := os.Stat(genSQLFile); os.IsNotExist(err) {
		return fmt.Errorf("SQL 文件不存在: %s", genSQLFile)
	}

	fmt.Printf("📄 正在解析 SQL 文件: %s\n\n", genSQLFile)

	generator := gen.NewSQLGenerator(gen.Config{
		SQLFile: genSQLFile,
		Output:  genOutput,
	})

	if err := generator.Generate(); err != nil {
		return fmt.Errorf("生成代码失败: %w", err)
	}

	fmt.Println("\n✅ 代码生成完成！")

	return nil
}

// selectTablesInteractive 交互式选择表
func selectTablesInteractive(dsn string) ([]string, error) {
	// 连接数据库，获取所有表
	fmt.Println("🔍 正在读取数据库表列表...")

	tables, err := gen.ListTables(dsn)
	if err != nil {
		return nil, fmt.Errorf("读取表列表失败: %w", err)
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("数据库中没有找到表")
	}

	fmt.Printf("\n📋 发现以下表（共 %d 张）：\n\n", len(tables))

	// 显示表列表
	for i, table := range tables {
		comment := table.Comment
		if comment == "" {
			comment = "-"
		}
		fmt.Printf("   [%2d] %-20s (%s) %3d 字段  %2d 索引\n",
			i+1, table.Name, comment, table.FieldsCount, table.IndexesCount)
	}

	fmt.Println("\n📝 请选择要生成的表：")
	fmt.Println("   方式：")
	fmt.Println("   - 输入序号（逗号分隔）: 1,2,3")
	fmt.Println("   - 输入范围: 1-5")
	fmt.Println("   - 输入通配符: user*")
	fmt.Println("   - 输入 all 生成所有表")
	fmt.Print("\n👉 您的选择: ")

	var input string
	fmt.Scanln(&input)

	// 解析输入
	return parseTableSelection(input, tables)
}

// parseTableSelection 解析表选择
func parseTableSelection(input string, tables []gen.TableInfo) ([]string, error) {
	input = strings.TrimSpace(input)

	// 输入 "all"
	if strings.ToLower(input) == "all" {
		var names []string
		for _, t := range tables {
			names = append(names, t.Name)
		}
		return names, nil
	}

	// 检查是否包含通配符
	if strings.Contains(input, "*") {
		return filterTablesByWildcard(input, tables)
	}

	// 解析序号或表名
	return parseTableInput(input, tables)
}

// parseTableInput 解析表输入（序号或表名）
func parseTableInput(input string, tables []gen.TableInfo) ([]string, error) {
	var selected []string
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// 检查是否是范围 (1-5)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) == 2 {
				start := parseInt(rangeParts[0])
				end := parseInt(rangeParts[1])
				for i := start; i <= end; i++ {
					if i > 0 && i <= len(tables) {
						selected = append(selected, tables[i-1].Name)
					}
				}
				continue
			}
		}

		// 尝试解析为数字（序号）
		index := parseInt(part)
		if index > 0 && index <= len(tables) {
			selected = append(selected, tables[index-1].Name)
		} else {
			// 作为表名处理
			selected = append(selected, part)
		}
	}

	return selected, nil
}

// filterTablesByWildcard 使用通配符过滤表
func filterTablesByWildcard(pattern string, tables []gen.TableInfo) ([]string, error) {
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	var selected []string

	for _, t := range tables {
		matched, err := regexp.MatchString(pattern, t.Name)
		if err != nil {
			return nil, err
		}
		if matched {
			selected = append(selected, t.Name)
		}
	}

	return selected, nil
}

// 辅助函数
func maskDSN(dsn string) string {
	// 简单的密码遮蔽
	if strings.Contains(dsn, ":") && strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "@")
		if len(parts) >= 2 {
			userPass := strings.Split(parts[0], ":")
			if len(userPass) >= 2 {
				return userPass[0] + ":****@" + parts[1]
			}
		}
	}
	return dsn
}

func parseTables(tables string) []string {
	if tables == "" {
		return []string{}
	}
	parts := strings.Split(tables, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func loadTablesFromConfig(filename string) ([]string, error) {
	// TODO: 实现配置文件解析
	return []string{}, fmt.Errorf("配置文件功能尚未实现")
}
