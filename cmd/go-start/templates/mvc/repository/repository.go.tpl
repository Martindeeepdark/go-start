package repository

import (
	"{{.Module}}/pkg/database"
)

// Repository 仓储层聚合
type Repository struct {
	User UserRepository
}

// New 创建新的仓储实例
func New(db *database.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db.GormDB()),
	}
}

