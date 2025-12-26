package gen

import (
	"database/sql"
	"fmt"
	"strings"

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
	// TODO: 实现代码生成
	return nil
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

// Generate 生成代码
func (g *SQLGenerator) Generate() error {
	// TODO: 实现代码生成
	return nil
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
