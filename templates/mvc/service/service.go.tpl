package service

import (
	"github.com/{{.Module}}/internal/repository"
	"github.com/{{.Module}}/pkg/cache"
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
