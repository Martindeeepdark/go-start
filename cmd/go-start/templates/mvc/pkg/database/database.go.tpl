package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config represents database configuration
type Config struct {
	Driver          string `yaml:"driver" mapstructure:"driver"`
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Database        string `yaml:"database" mapstructure:"database"`
	Username        string `yaml:"username" mapstructure:"username"`
	Password        string `yaml:"password" mapstructure:"password"`
	Charset         string `yaml:"charset" mapstructure:"charset"`
	ParseTime       bool   `yaml:"parse_time" mapstructure:"parse_time"`
	Loc             string `yaml:"loc" mapstructure:"loc"`
	MaxIdleConns    int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	LogLevel        string `yaml:"log_level" mapstructure:"log_level"`
}

// DB represents the GORM database wrapper
type DB struct {
	db *gorm.DB
}

// New creates a new database instance
func New(cfg *Config) (*DB, error) {
	dsn := buildDSN(cfg)

	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = newMySQLDriver(dsn)
	case "postgres":
		dialector = newPostgresDriver(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	logLevel := parseLogLevel(cfg.LogLevel)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db: db}, nil
}

// GormDB returns the underlying gorm.DB instance
func (d *DB) GormDB() *gorm.DB {
	return d.db
}

// Close closes the database connection
func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func buildDSN(cfg *Config) string {
	switch cfg.Driver {
	case "mysql":
		charset := cfg.Charset
		if charset == "" {
			charset = "utf8mb4"
		}
		loc := cfg.Loc
		if loc == "" {
			loc = "Local"
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, charset, cfg.ParseTime, loc)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	default:
		return ""
	}
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
