package user

import (
	"context"
	"fmt"

	"{{.Module}}/internal/domain/user"
)

// UserService 用户应用服务
type UserService struct {
	userRepo user.UserRepository
}

// NewUserService 创建用户应用服务
func NewUserService(userRepo user.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户用例
func (s *UserService) CreateUser(ctx context.Context, username, email, password string, age int) error {
	// 1. 创建领域对象
	u, err := user.NewUser(username, email, password, age)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	// 2. 保存到仓储
	if err := s.userRepo.Save(ctx, u); err != nil {
		return fmt.Errorf("保存用户失败: %w", err)
	}

	return nil
}

// UpdateUserEmail 更新用户邮箱用例
func (s *UserService) UpdateUserEmail(ctx context.Context, userID uint, newEmail string) error {
	// 1. 查询聚合根
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 执行业务操作
	if err := u.ChangeEmail(newEmail); err != nil {
		return err
	}

	// 3. 保存更改
	if err := s.userRepo.Save(ctx, u); err != nil {
		return fmt.Errorf("保存用户失败: %w", err)
	}

	return nil
}
