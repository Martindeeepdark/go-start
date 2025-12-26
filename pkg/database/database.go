package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TxOptions represents transaction options
type TxOptions struct {
	Isolation int
	ReadOnly  bool
}

// Stats represents database statistics
type Stats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUse              int
	Idle               int
}

// Config represents database configuration
type Config struct {
	Driver          string `yaml:"driver" mapstructure:"driver"` // mysql, postgres
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
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"` // seconds
	LogLevel        string `yaml:"log_level" mapstructure:"log_level"`                 // silent, error, warn, info
}

// DB represents the GORM database wrapper
type DB struct {
	db *gorm.DB
}

// New creates a new database instance
func New(cfg *Config) (*DB, error) {
	var dialector gorm.Dialector

	dsn := buildDSN(cfg)
	switch cfg.Driver {
	case "mysql":
		dialector = NewMySQLDriver(dsn)
	case "postgres":
		dialector = NewPostgresDriver(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Parse log level
	logLevel := parseLogLevel(cfg.LogLevel)

	// Create GORM instance
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db: db}, nil
}

// buildDSN builds the data source name
func buildDSN(cfg *Config) string {
	switch cfg.Driver {
	case "mysql":
		return buildMySQLDSN(cfg)
	case "postgres":
		return buildPostgresDSN(cfg)
	default:
		return ""
	}
}

// buildMySQLDSN builds MySQL DSN
func buildMySQLDSN(cfg *Config) string {
	charset := cfg.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	parseTime := "True"
	if !cfg.ParseTime {
		parseTime = "False"
	}

	loc := cfg.Loc
	if loc == "" {
		loc = "Local"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, charset, parseTime, loc)
}

// buildPostgresDSN builds PostgreSQL DSN
func buildPostgresDSN(cfg *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
}

// parseLogLevel parses log level string
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
		return logger.Info
	}
}

// DB returns the underlying GORM DB instance
func (d *DB) DB() *gorm.DB {
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

// BeginTx starts a transaction
func (d *DB) BeginTx(ctx context.Context, opts *TxOptions) (*Transaction, error) {
	var sqlOpts *sql.TxOptions
	if opts != nil {
		sqlOpts = &sql.TxOptions{
			Isolation: sql.IsolationLevel(opts.Isolation),
			ReadOnly:  opts.ReadOnly,
		}
	}

	tx := d.db.WithContext(ctx).Begin(sqlOpts)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &Transaction{tx: tx}, nil
}

// Transaction represents a GORM transaction
type Transaction struct {
	tx *gorm.DB
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	return t.tx.Commit().Error
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	return t.tx.Rollback().Error
}

// Tx returns the underlying GORM transaction DB
func (t *Transaction) Tx() *gorm.DB {
	return t.tx
}

// Ping pings the database
func (d *DB) Ping(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Stats returns database statistics
func (d *DB) Stats() Stats {
	sqlDB, err := d.db.DB()
	if err != nil {
		return Stats{}
	}

	stats := sqlDB.Stats()
	return Stats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUse:              stats.InUse,
		Idle:               stats.Idle,
	}
}
