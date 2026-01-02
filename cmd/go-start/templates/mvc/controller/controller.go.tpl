package controller

import (
	"{{.Module}}/internal/service"
)

// Controller represents the controller layer
type Controller struct {
	User *UserController
}

// New creates a new controller instance
func New(service *service.Service) *Controller {
	return &Controller{
		User: NewUserController(service),
	}
}
