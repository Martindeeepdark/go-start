package persistence

import (
	"github.com/{{.Module}}/internal/domain/user"
	"github.com/{{.Module}}/internal/infrastructure/persistence/user"
	"gorm.io/gorm"
)

// Repositories 仓储集合
type Repositories struct {
	User user.UserRepository
	// TODO: 添加其他聚合的仓储
}

// NewRepositories 初始化所有仓储
func NewRepositories(db *gorm.DB, cache interface{}) *Repositories {
	return &Repositories{
		User: user.NewUserRepositoryImpl(db),
		// TODO: 初始化其他仓储
	}
}
