package gen

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// TableInfo 表信息
type TableInfo struct {
	Name         string // 表名
	Comment      string // 表注释
	FieldsCount  int    // 字段数量
	IndexesCount int    // 索引数量
}

// FieldInfo 字段信息
type FieldInfo struct {
	Name          string // 字段名
	Type          string // 数据库类型
	GoType        string // Go 类型
	Nullable      bool   // 是否可为空
	PrimaryKey    bool   // 是否主键
	AutoIncrement bool   // 是否自增
	DefaultValue  string // 默认值
	Comment       string // 注释
}

// IndexInfo 索引信息
type IndexInfo struct {
	Name    string   // 索引名
	Columns []string // 索引列
	Unique  bool     // 是否唯一索引
	Primary bool     // 是否主键索引
}

// Config 生成器配置
type Config struct {
	DSN     string   // 数据库连接字符串
	Tables  []string // 要生成的表名
	Output  string   // 输出目录
	SQLFile string   // SQL 文件路径（用于 SQL 生成器）
	Module  string   // Go 模块路径
}

// DatabaseGenerator 数据库代码生成器
type DatabaseGenerator struct {
	config Config
}

// NewDatabaseGenerator 创建数据库代码生成器
func NewDatabaseGenerator(config Config) *DatabaseGenerator {
	return &DatabaseGenerator{
		config: config,
	}
}

// Generate 生成代码
func (g *DatabaseGenerator) Generate() error {
	fmt.Println("🔧 正在初始化 GORM Gen...")

	// 0. 检查 Go 版本
	if err := checkGoVersion(); err != nil {
		return err
	}

	// 1. 连接数据库
	db, err := connectGORMDB(g.config.DSN)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 2. 创建 GORM Gen 生成器
	generator := gen.NewGenerator(gen.Config{
		OutPath:       filepath.Join(g.config.Output, "internal/dal/query"),
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
		FieldSignable: false,
		// FieldWithIndexTag: false,
		// FieldWithTypeTag: true,
	})

	generator.UseDB(db)

	// 3. 为每个表生成模型
	fmt.Println("📦 正在读取表结构...")
	var models []interface{}

	for _, tableName := range g.config.Tables {
		fmt.Printf("  📋 处理表: %s\n", tableName)

		// 🔧 检查表是否存在
		if !g.tableExists(db, tableName) {
			return fmt.Errorf("❌ 表 '%s' 不存在\n\n💡 帮助:\n%s\n📋 数据库中的表:\n%s",
				tableName,
				g.getHint(),
				g.listTables(db))
		}

		// 使用 GORM Gen 自动生成模型
		model := generator.GenerateModel(tableName)
		models = append(models, model)
	}

	// 4. 执行 GORM Gen 生成
	fmt.Println("🚀 正在生成 GORM 查询代码...")
	generator.ApplyBasic(models...)
	generator.Execute()

	fmt.Println("✅ GORM Gen 代码生成完成！")
	fmt.Printf("   生成位置: %s\n", filepath.Join(g.config.Output, "internal/dal/query"))

	// 5. 生成 Repository 层
	fmt.Println("\n📦 正在生成 Repository 层...")
	if err := g.generateRepositoryLayer(); err != nil {
		return fmt.Errorf("生成 Repository 层失败: %w", err)
	}

	// 6. 生成 Service 层
	fmt.Println("\n📦 正在生成 Service 层...")
	if err := g.generateServiceLayer(); err != nil {
		return fmt.Errorf("生成 Service 层失败: %w", err)
	}

	// 7. 生成 Controller 层
	fmt.Println("\n📦 正在生成 Controller 层...")
	if err := g.generateControllerLayer(); err != nil {
		return fmt.Errorf("生成 Controller 层失败: %w", err)
	}

	// 8. 生成 go.mod
	fmt.Println("\n📦 正在生成 go.mod...")
	if err := g.GenerateGoMod(); err != nil {
		return fmt.Errorf("生成 go.mod 失败: %w", err)
	}

	// 9. 生成支持包
	fmt.Println("\n📦 正在生成支持包...")
	if err := g.GenerateSupportPackages(); err != nil {
		return fmt.Errorf("生成支持包失败: %w", err)
	}

	// 10. 生成 main.go 和配置文件
	fmt.Println("\n📦 正在生成 main.go...")
	if err := g.GenerateMainGo(); err != nil {
		return fmt.Errorf("生成 main.go 失败: %w", err)
	}

	// 11. 生成路由注册
	fmt.Println("\n📦 正在生成路由注册...")
	if err := g.GenerateRoutes(getTablesFromNames(g.config.Tables), getModulePath(g.config.Module)); err != nil {
		return fmt.Errorf("生成路由失败: %w", err)
	}

	fmt.Println("\n✅ 所有代码生成完成！")

	return nil
}

// getTablesFromNames 从表名列表创建 TableInfo 列表
func getTablesFromNames(names []string) []TableInfo {
	var tables []TableInfo
	for _, name := range names {
		tables = append(tables, TableInfo{Name: name})
	}
	return tables
}

// connectGORMDB 使用 GORM 连接数据库
func connectGORMDB(dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// 根据DSN判断数据库类型
	if strings.Contains(dsn, "@tcp(") || strings.Contains(dsn, "@unix(") {
		// MySQL
		dialector = mysql.Open(dsn)
	} else if strings.HasPrefix(dsn, "host=") || strings.HasPrefix(dsn, "postgres://") {
		// PostgreSQL
		dialector = postgres.Open(dsn)
	} else {
		return nil, fmt.Errorf("不支持的数据库类型，DSN 格式无法识别")
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	return db, nil
}

// tableExists 检查表是否存在
func (g *DatabaseGenerator) tableExists(db *gorm.DB, tableName string) bool {
	var count int64
	// 查询当前数据库中的表
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Scan(&count)
	return count > 0
}

// getHint 获取提示信息
func (g *DatabaseGenerator) getHint() string {
	return `• 检查表名是否拼写正确
• 使用 --interactive 交互式选择表
• 使用 --tables="*" 生成所有表
• 使用 --tables="prefix*" 生成指定前缀的表`
}

// listTables 列出数据库中的所有表
func (g *DatabaseGenerator) listTables(db *gorm.DB) string {
	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() ORDER BY table_name").Scan(&tables)

	if len(tables) == 0 {
		return "  (空)"
	}

	var result string
	for _, table := range tables {
		result += fmt.Sprintf("  • %s\n", table)
	}
	return result
}

// generateRepositoryLayer 生成 Repository 层
func (g *DatabaseGenerator) generateRepositoryLayer() error {
	// TODO: 获取模块路径和索引信息
	// 这里简化处理，实际应该从配置或读取生成的代码获取

	for _, tableName := range g.config.Tables {
		// 简单的表名转模型名
		modelName := toModelName(tableName)

		// 获取表的索引信息
		schema, err := GetTableSchema(g.config.DSN, tableName)
		if err != nil {
			fmt.Printf("  ⚠️  无法获取 %s 的索引信息，跳过\n", tableName)
			continue
		}

		// 提取索引字段（去重）
		indexFieldSet := make(map[string]bool)
		var indexFields []string
		for _, idx := range schema.Indexes {
			if !idx.Primary && len(idx.Columns) == 1 {
				// 只处理单列非主键索引
				field := toCamelCase(idx.Columns[0])
				if !indexFieldSet[field] {
					indexFieldSet[field] = true
					indexFields = append(indexFields, field)
				}
			}
		}

		// 配置 Repository 生成
		config := RepositoryConfig{
			TableName:   tableName,
			ModelName:   modelName,
			PackageName: "repository",
			ModulePath:  getModulePath(g.config.Module),
			Indexes:     indexFields,
		}

		if err := g.GenerateRepository(TableInfo{Name: tableName}, config); err != nil {
			return err
		}
	}

	return nil
}

// generateServiceLayer 生成 Service 层
func (g *DatabaseGenerator) generateServiceLayer() error {
	for _, tableName := range g.config.Tables {
		modelName := toModelName(tableName)

		config := ServiceConfig{
			TableName:   tableName,
			ModelName:   modelName,
			PackageName: "service",
			ModulePath:  getModulePath(g.config.Module),
		}

		if err := g.GenerateService(TableInfo{Name: tableName}, config); err != nil {
			return err
		}
	}

	return nil
}

// generateControllerLayer 生成 Controller 层
func (g *DatabaseGenerator) generateControllerLayer() error {
	for _, tableName := range g.config.Tables {
		// 简单的表名转模型名
		modelName := toModelName(tableName)

		// 配置 Controller 生成
		config := ControllerConfig{
			TableName:   tableName,
			ModelName:   modelName,
			PackageName: "controller",
			ModulePath:  getModulePath(g.config.Module),
		}

		if err := g.GenerateController(TableInfo{Name: tableName}, config); err != nil {
			return err
		}
	}

	return nil
}

// getModulePath 获取模块路径，如果为空则自动检测
func getModulePath(modulePath string) string {
	if modulePath == "" {
		// 自动检测模块路径
		if detected := detectModulePath(); detected != "" {
			return detected
		}
		// 如果检测失败，返回空字符串，让调用者处理
		return ""
	}
	return modulePath
}

// detectModulePath 从 go.mod 文件检测模块路径
func detectModulePath() string {
	// 尝试读取当前目录的 go.mod
	if modulePath := extractModulePath("go.mod"); modulePath != "" {
		return modulePath
	}
	// 尝试读取父目录的 go.mod
	// TODO: 可以向上递归查找
	return ""
}

// extractModulePath 从 go.mod 文件提取模块路径
func extractModulePath(goModPath string) string {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

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

// toModelName 表名转模型名
func toModelName(tableName string) string {
	// users -> Users
	// user_profiles -> UserProfile
	parts := strings.Split(tableName, "_")
	for i, part := range parts {
		if i > 0 || len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// toCamelCase 转换为驼峰命名
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// SQLGenerator SQL文件代码生成器
type SQLGenerator struct {
	config Config
}

// NewSQLGenerator 创建SQL文件代码生成器
func NewSQLGenerator(config Config) *SQLGenerator {
	return &SQLGenerator{
		config: config,
	}
}

// Generate 从 SQL 文件生成代码
func (g *SQLGenerator) Generate() error {
	fmt.Println("📄 正在解析 SQL 文件...")

	// 1. 读取 SQL 文件
	_, err := os.ReadFile(g.config.SQLFile)
	if err != nil {
		return fmt.Errorf("读取 SQL 文件失败: %w", err)
	}

	// 2. 创建临时数据库并导入 SQL
	// 注意：这需要一个临时数据库
	// 简化起见，我们提示用户使用 DatabaseGenerator

	fmt.Println("⚠️  SQL 文件生成功能建议使用以下方式：")
	fmt.Println()
	fmt.Println("方式一：直接使用数据库生成（推荐）")
	fmt.Println("  1. 创建数据库并导入 SQL 文件：")
	fmt.Println("     mysql -u root -p < schema.sql")
	fmt.Println()
	fmt.Println("  2. 使用数据库生成命令：")
	fmt.Printf("     go-start gen db --dsn=\"root:pass@tcp(localhost:3306)/dbname\" --tables=your_tables\n")
	fmt.Println()
	fmt.Println("方式二：使用交互式选择")
	fmt.Println("  go-start gen db --dsn=\"...\" --interactive")
	fmt.Println()
	fmt.Println("为什么推荐使用数据库连接？")
	fmt.Println("  ✅ 可以准确读取表结构")
	fmt.Println("  ✅ 可以获取索引信息")
	fmt.Println("  ✅ 可以生成基于索引的查询方法")
	fmt.Println("  ✅ 支持交互式表选择")

	return fmt.Errorf("请使用 DatabaseGenerator 从数据库生成代码")
}

// ListTables 列出数据库中的所有表
func ListTables(dsn string) ([]TableInfo, error) {
	return GetTables(dsn)
}

// GetTables 获取数据库中的所有表（带字段和索引统计）
func GetTables(dsn string) ([]TableInfo, error) {
	// 解析 DSN 获取数据库类型
	dbType, dbName, err := parseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("解析 DSN 失败: %w", err)
	}

	// 连接数据库
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 获取表列表
	tables, err := getTables(db, dbType, dbName)
	if err != nil {
		return nil, fmt.Errorf("获取表列表失败: %w", err)
	}

	// 获取每张表的字段和索引数量
	for i := range tables {
		fields, _ := getFields(db, dbType, tables[i].Name)
		indexes, _ := getIndexes(db, dbType, tables[i].Name)
		tables[i].FieldsCount = len(fields)
		tables[i].IndexesCount = len(indexes)
	}

	return tables, nil
}

// GetTableSchema 获取表的详细结构
func GetTableSchema(dsn, tableName string) (*DetailedTableInfo, error) {
	dbType, _, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 获取字段信息
	fields, err := getFields(db, dbType, tableName)
	if err != nil {
		return nil, err
	}

	// 获取索引信息
	indexes, err := getIndexes(db, dbType, tableName)
	if err != nil {
		return nil, err
	}

	return &DetailedTableInfo{
		Name:    tableName,
		Fields:  fields,
		Indexes: indexes,
	}, nil
}

// DetailedTableInfo 详细的表信息
type DetailedTableInfo struct {
	Name    string      // 表名
	Fields  []FieldInfo // 字段列表
	Indexes []IndexInfo // 索引列表
}

// parseDSN, getTables, getFields, getIndexes, mapToGoType 等函数...
// (从 database.go 复制过来)

// parseDSN 解析 DSN，返回数据库类型、数据库名
func parseDSN(dsn string) (dbType, dbName string, err error) {
	if strings.Contains(dsn, "@tcp(") || strings.Contains(dsn, "@unix(") {
		// MySQL 格式: user:pass@tcp(host:port)/dbname
		parts := strings.Split(dsn, "/")
		if len(parts) >= 2 {
			dbType = "mysql"
			dbName = parts[len(parts)-1]
			if idx := strings.Index(dbName, "?"); idx > 0 {
				dbName = dbName[:idx]
			}
			return
		}
	} else if strings.HasPrefix(dsn, "host=") || strings.HasPrefix(dsn, "postgres://") {
		// PostgreSQL
		dbType = "postgres"
		if strings.Contains(dsn, "dbname=") {
			parts := strings.Split(dsn, "dbname=")
			if len(parts) > 1 {
				dbName = parts[1]
				if idx := strings.Index(dbName, " "); idx > 0 {
					dbName = dbName[:idx]
				}
			}
		}
		return
	}
	err = fmt.Errorf("不支持的 DSN 格式")
	return
}

// getTables 获取表列表
func getTables(db *sql.DB, dbType, dbName string) ([]TableInfo, error) {
	var query string
	var args []interface{}

	if dbType == "mysql" {
		query = `SELECT TABLE_NAME, TABLE_COMMENT
			FROM INFORMATION_SCHEMA.TABLES
			WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE'
			ORDER BY TABLE_NAME`
		args = []interface{}{dbName}
	} else if dbType == "postgres" {
		query = `SELECT table_name, obj_description((table_schema||'.'||table_name)::regclass, 'oid') as comment
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
			ORDER BY table_name`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var name, comment sql.NullString
		if err := rows.Scan(&name, &comment); err != nil {
			return nil, err
		}

		tables = append(tables, TableInfo{
			Name:    name.String,
			Comment: comment.String,
		})
	}

	return tables, nil
}

// getFields 获取字段列表
func getFields(db *sql.DB, dbType, tableName string) ([]FieldInfo, error) {
	var query string

	if dbType == "mysql" {
		query = `SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_KEY,
			COLUMN_DEFAULT, EXTRA, COLUMN_COMMENT
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
			ORDER BY ORDINAL_POSITION`
	} else if dbType == "postgres" {
		query = `SELECT column_name, data_type, is_nullable,
			column_default, column_comment
			FROM information_schema.columns
			WHERE table_name = $1
			ORDER BY ordinal_position`
	}

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []FieldInfo
	for rows.Next() {
		var field FieldInfo
		var nullable sql.NullString
		var columnKey, extra, defaultValue, comment sql.NullString

		if dbType == "mysql" {
			if err := rows.Scan(&field.Name, &field.Type, &nullable,
				&columnKey, &defaultValue, &extra, &comment); err != nil {
				return nil, err
			}
			field.PrimaryKey = columnKey.String == "PRI"
			field.AutoIncrement = strings.Contains(extra.String, "auto_increment")
			field.Comment = comment.String
			field.DefaultValue = defaultValue.String
		} else {
			if err := rows.Scan(&field.Name, &field.Type, &nullable,
				&defaultValue, &comment); err != nil {
				return nil, err
			}
			field.Comment = comment.String
			field.DefaultValue = defaultValue.String
		}

		field.Nullable = nullable.String == "YES"
		field.GoType = mapToGoType(field.Type)

		fields = append(fields, field)
	}

	return fields, nil
}

// getIndexes 获取索引列表
func getIndexes(db *sql.DB, dbType, tableName string) ([]IndexInfo, error) {
	var query string

	if dbType == "mysql" {
		query = `SELECT INDEX_NAME, GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) as COLUMNS, NON_UNIQUE
			FROM INFORMATION_SCHEMA.STATISTICS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
			GROUP BY INDEX_NAME, NON_UNIQUE ORDER BY INDEX_NAME`
	} else if dbType == "postgres" {
		// PostgreSQL 索引查询较复杂，这里简化处理
		return []IndexInfo{}, nil
	}

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []IndexInfo
	for rows.Next() {
		var name, columns sql.NullString
		var nonUnique sql.NullInt64

		if err := rows.Scan(&name, &columns, &nonUnique); err != nil {
			return nil, err
		}

		indexInfo := IndexInfo{
			Name:    name.String,
			Columns: strings.Split(columns.String, ","),
			Unique:  nonUnique.Int64 == 0,
			Primary: name.String == "PRIMARY",
		}

		indexes = append(indexes, indexInfo)
	}

	return indexes, nil
}

// mapToGoType 数据库类型映射到 Go 类型
func mapToGoType(dbType string) string {
	typeMap := map[string]string{
		"varchar": "string", "char": "string", "text": "string",
		"int": "int", "tinyint": "int", "smallint": "int",
		"bigint": "int64", "float": "float64", "double": "float64",
		"datetime": "time.Time", "timestamp": "time.Time",
		"date": "time.Time", "bool": "bool", "boolean": "bool",
	}

	if goType, ok := typeMap[strings.ToLower(dbType)]; ok {
		return goType
	}
	return "string"
}
