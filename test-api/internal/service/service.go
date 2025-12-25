package service

import (
	"github.com/github.com/yourname/test-api/internal/repository"
	"github.com/github.com/yourname/test-api/pkg/cache"
)

// Service represents the service layer
type Service struct {
	User *UserService
}

// New creates a new service instance
func New(repo *repository.Repository, cache *cache.Cache) *Service {
	return &Service{
		User: NewUserService(repo.User, cache),
	}
}
