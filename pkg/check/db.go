package check

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig holds database configuration
type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// TestDatabaseConnection tests if the database connection is working
func TestDatabaseConnection(config *DBConfig) error {
	var dsn string
	var dialector gorm.Dialector

	switch strings.ToLower(config.Driver) {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)
		dialector = mysql.Open(dsn)

	case "postgresql", "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host,
			config.Port,
			config.Username,
			config.Password,
			config.Database,
		)
		dialector = postgres.Open(dsn)

	default:
		return fmt.Errorf("不支持的数据库类型: %s (支持: mysql, postgresql)", config.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	defer sqlDB.Close()

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库 Ping 失败: %w", err)
	}

	return nil
}

// PrintDatabaseTestResult prints the database connection test result
func PrintDatabaseTestResult(err error) {
	if err != nil {
		fmt.Println("❌ 数据库连接测试失败")
		fmt.Println("   错误:", err.Error())
		fmt.Println()
		fmt.Println("💡 请检查:")
		fmt.Println("   1. 数据库服务是否启动")
		fmt.Println("   2. 配置文件中的连接信息是否正确")
		fmt.Println("   3. 数据库是否已创建")
		fmt.Println("   4. 用户权限是否正确")
		fmt.Println()
	} else {
		fmt.Println("✅ 数据库连接测试成功")
		fmt.Println()
	}
}
