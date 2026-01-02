package repository

import (
	"{{.Module}}/pkg/database"
)

// Repository represents the repository layer
type Repository struct {
	User *UserRepository
}

// New creates a new repository instance
func New(db *database.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db.DB()),
	}
}
