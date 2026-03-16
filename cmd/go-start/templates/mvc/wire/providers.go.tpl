package providers

import (
	"{{.Module}}/config"
	"{{.Module}}/internal/controller"
	"{{.Module}}/internal/repository"
	"{{.Module}}/internal/service"
	"{{.Module}}/pkg/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

func ProvideLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func ProvideDatabase(cfg *config.Config) (*database.DB, error) {
	return database.New(&cfg.Database)
}

func ProvideRepository(db *database.DB) repository.UserRepository {
	return repository.NewUserRepository(db.DB())
}

func ProvideService(repo repository.UserRepository) *service.UserService {
	return service.NewUserService(repo)
}

func ProvideController(svc *service.UserService) *controller.UserController {
	return controller.NewUserController(svc)
}

func ProvideRouter(ctrl *controller.UserController) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	RegisterRoutes(r, ctrl)
	return r
}

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
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
