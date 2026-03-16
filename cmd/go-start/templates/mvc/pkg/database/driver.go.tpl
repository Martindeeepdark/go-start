package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMySQLDriver(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

func newPostgresDriver(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
