package user

import "context"

// UserRepository 用户仓储接口
//
// 职责说明：
//   - 定义用户的持久化抽象
//   - 隐藏基础设施细节
//   - 由 Infrastructure 层实现
type UserRepository interface {
	// Save 保存用户
	Save(ctx context.Context, user *User) error

	// FindByID 根据 ID 查找用户
	FindByID(ctx context.Context, id uint) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindAll 查找所有用户
	FindAll(ctx context.Context) ([]*User, error)

	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
}
