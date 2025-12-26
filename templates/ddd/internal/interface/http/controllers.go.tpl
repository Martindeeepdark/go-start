package http

import (
	"github.com/{{.Module}}/internal/application"
)

// Controllers 控制器集合
type Controllers struct {
	User *UserController
	// TODO: 添加其他控制器
}

// NewControllers 初始化所有控制器
func NewControllers(apps *application.Applications) *Controllers {
	return &Controllers{
		User: NewUserController(apps.User),
		// TODO: 初始化其他控制器
	}
}
