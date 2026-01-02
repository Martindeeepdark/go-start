package providers

import (
	"{{.Module}}/config"
	"{{.Module}}/internal/controller"
	"{{.Module}}/internal/repository"
	"{{.Module}}/internal/service"
	"{{.Module}}/pkg/cache"
	"{{.Module}}/pkg/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ============================================
// Provider 函数：提供依赖
// ============================================

// ProvideConfig 提供配置
func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

// ProvideLogger 提供日志
func ProvideLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

// ProvideDatabase 提供数据库
func ProvideDatabase(cfg *config.Config) (*database.DB, error) {
	return database.New(&cfg.Database)
}

// ProvideCache 提供 Redis（可选）
func ProvideCache(cfg *config.Config) (*cache.Cache, error) {
	return cache.New(&cfg.Redis)
}

// ProvideRepository 提供仓储层
//
// 说明：这里使用接口而不是具体实现
// 便于测试和切换实现
func ProvideRepository(db *database.DB) repository.UserRepository {
	return repository.NewUserRepository(db.DB())
}

// ProvideService 提供服务层
//
// 说明：Service 层依赖 Repository 接口
// 符合依赖倒置原则（DIP）
func ProvideService(repo repository.UserRepository, cache *cache.Cache) *service.UserService {
	return service.NewUserService(repo, cache)
}

// ProvideController 提供控制器层
func ProvideController(svc *service.UserService) *controller.UserController {
	return controller.NewUserController(svc)
}

// ProvideRouter 提供路由
//
// 说明：将所有依赖注入到路由中
// 这里演示了完整的依赖链：DB → Repo → Service → Controller → Router
func ProvideRouter(ctrl *controller.UserController, logger *zap.Logger) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// 注册路由
	RegisterRoutes(r, ctrl)

	return r
}

// ============================================
// 路由注册
// ============================================

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, ctrl *controller.UserController) {
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", ctrl.Create)
			users.GET(":id", ctrl.GetByID)
			users.PUT(":id", ctrl.Update)
			users.DELETE(":id", ctrl.Delete)
			users.GET("", ctrl.List)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
