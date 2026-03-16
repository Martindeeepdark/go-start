package service

import (
	"{{.Module}}/internal/repository"
)

// Service represents the service layer
type Service struct {
	User *UserService
}

// New creates a new service instance
func New(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo.User),
	}
}
