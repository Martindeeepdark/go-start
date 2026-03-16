package controller

import (
	"github.com/gin-gonic/gin"
	"{{.Module}}/internal/service"
)

// Controller represents the controller layer
type Controller struct {
	User *UserController
}

// New creates a new controller instance
func New(svc *service.Service) *Controller {
	return &Controller{
		User: NewUserController(svc),
	}
}

// RegisterRoutes registers all routes
func (c *Controller) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", c.User.Create)
			users.GET("/:id", c.User.Get)
			users.PUT("/:id", c.User.Update)
			users.DELETE("/:id", c.User.Delete)
			users.GET("", c.User.List)
		}
	}
}
