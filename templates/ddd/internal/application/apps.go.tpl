package application

import (
	"github.com/{{.Module}}/internal/application/user"
	"github.com/{{.Module}}/internal/infrastructure/persistence"
)

// Applications 应用服务集合
type Applications struct {
	User *user.UserService
	// 在此添加其他应用服务 (如: ArticleService, CommentService 等)
}

// NewApplications 初始化所有应用服务
func NewApplications(repos *persistence.Repositories) *Applications {
	return &Applications{
		User: user.NewUserService(repos.User),
		// 在此初始化其他应用服务
	}
}
