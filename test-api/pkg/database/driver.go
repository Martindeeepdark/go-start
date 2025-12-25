package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMySQLDriver creates a MySQL driver
func NewMySQLDriver(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

// NewPostgresDriver creates a PostgreSQL driver
func NewPostgresDriver(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
