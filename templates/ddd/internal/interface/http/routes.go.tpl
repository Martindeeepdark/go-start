package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, controllers *Controllers) {
	api := r.Group("/api")
	{
		// User routes
		userGroup := api.Group("/users")
		{
			userGroup.POST("", controllers.User.CreateUser)
			userGroup.PUT("/:id/email", controllers.User.UpdateEmail)
			// 在此添加更多用户相关路由
		}

		// 在此添加其他聚合的路由 (如: /articles, /comments 等)
	}
}
