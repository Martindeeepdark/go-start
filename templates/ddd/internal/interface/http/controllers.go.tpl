package http

import (
	"github.com/{{.Module}}/internal/application"
)

// Controllers 控制器集合
type Controllers struct {
	User *UserController
	// 在此添加其他控制器 (如: ArticleController, CommentController 等)
}

// NewControllers 初始化所有控制器
func NewControllers(apps *application.Applications) *Controllers {
	return &Controllers{
		User: NewUserController(apps.User),
		// 在此初始化其他控制器
	}
}
