package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"github.com/{{.Module}}/internal/domain/user"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepositoryImpl 创建用户仓储实现
func NewUserRepositoryImpl(db *gorm.DB) user.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save 保存用户
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
	// TODO: 领域模型与数据模型转换
	// dataModel := r.toDataModel(u)
	// return r.db.WithContext(ctx).Save(dataModel).Error
	return fmt.Errorf("not implemented")
}

// FindByID 根据 ID 查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*user.User, error) {
	// TODO: 查询并转换为领域模型
	// var dataModel UserDataModel
	// err := r.db.WithContext(ctx).First(&dataModel, id).Error
	// if err != nil {
	// 	return nil, err
	// }
	// return r.toDomainModel(&dataModel), nil
	return nil, fmt.Errorf("not implemented")
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, fmt.Errorf("not implemented")
}

// FindAll 查找所有用户
func (r *UserRepositoryImpl) FindAll(ctx context.Context) ([]*user.User, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return fmt.Errorf("not implemented")
}
