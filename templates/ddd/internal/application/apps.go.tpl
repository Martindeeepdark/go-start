package application

import (
	"github.com/{{.Module}}/internal/application/user"
	"github.com/{{.Module}}/internal/infrastructure/persistence"
)

// Applications 应用服务集合
type Applications struct {
	User *user.UserService
	// TODO: 添加其他应用服务
}

// NewApplications 初始化所有应用服务
func NewApplications(repos *persistence.Repositories) *Applications {
	return &Applications{
		User: user.NewUserService(repos.User),
		// TODO: 初始化其他应用服务
	}
}
