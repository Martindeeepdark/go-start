package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user model
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email     string         `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Status    int            `gorm:"default:1;comment:1-active,0-inactive" json:"status"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
