package service

import (
	"context"
	"errors"
	"fmt"

	"{{.Module}}/internal/model"
	"{{.Module}}/internal/repository"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrUserExists   = errors.New("用户已存在")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	existing, err := s.repo.FindByUsername(ctx, user.Username)
	if err == nil && existing != nil {
		return ErrUserExists
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	if _, err := s.repo.FindByID(ctx, user.ID); err != nil {
		return ErrUserNotFound
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return ErrUserNotFound
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}
	return nil
}

func (s *UserService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}
