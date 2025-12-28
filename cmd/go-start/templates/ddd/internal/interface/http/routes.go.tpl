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
			// TODO: 添加其他路由
		}

		// TODO: 添加其他聚合的路由
	}
}
