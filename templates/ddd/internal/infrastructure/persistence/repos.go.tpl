package persistence

import (
	"github.com/{{.Module}}/internal/domain/user"
	"github.com/{{.Module}}/internal/infrastructure/persistence/user"
	"gorm.io/gorm"
)

// Repositories 仓储集合
type Repositories struct {
	User user.UserRepository
	// 在此添加其他聚合的仓储 (如: Article, Comment 等)
}

// NewRepositories 初始化所有仓储
func NewRepositories(db *gorm.DB, cache interface{}) *Repositories {
	return &Repositories{
		User: user.NewUserRepositoryImpl(db),
		// 在此初始化其他仓储
	}
}
