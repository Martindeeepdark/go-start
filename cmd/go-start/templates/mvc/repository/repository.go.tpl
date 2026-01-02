package repository

import (
	"{{.Module}}/pkg/database"
	"gorm.io/gorm"
)

// Repository 仓储层聚合
type Repository struct {
	User UserRepository
}

// New 创建新的仓储实例
func New(db *database.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db.DB()),
	}
}

// 便捷方法：获取 GORM DB 实例（用于复杂查询）
func (r *Repository) DB() *gorm.DB {
	// 如果需要直接访问 DB，可以通过 repository 的某个字段获取
	// 或者返回 nil，强制使用 Repository 接口
	return nil
}
