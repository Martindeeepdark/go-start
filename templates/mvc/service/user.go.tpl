package service

import (
	"context"
	"errors"

	"github.com/{{.Module}}/internal/model"
	"github.com/{{.Module}}/internal/repository"
	"github.com/{{.Module}}/pkg/cache"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
)

// UserService represents the user service
type UserService struct {
	repo  *repository.UserRepository
	cache *cache.Cache
}

// NewUserService creates a new user service instance
func NewUserService(repo *repository.UserRepository, cache *cache.Cache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, user *model.User) error {
	// Check if user exists
	existing, err := s.repo.GetByUsername(ctx, user.Username)
	if err == nil && existing != nil {
		return ErrUserExists
	}

	return s.repo.Create(ctx, user)
}

// GetByID gets a user by ID
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	// Try to get from cache first
	// cacheKey := fmt.Sprintf("user:%d", id)
	// cached, err := s.cache.Get(ctx, cacheKey)
	// if err == nil {
	//     var user model.User
	//     if err := json.Unmarshal([]byte(cached), &user); err == nil {
	//         return &user, nil
	//     }
	// }

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Cache the user
	// if data, err := json.Marshal(user); err == nil {
	//     s.cache.Set(ctx, cacheKey, data, 10*time.Minute)
	// }

	return user, nil
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// List lists users with pagination
func (s *UserService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.List(ctx, page, pageSize)
}
