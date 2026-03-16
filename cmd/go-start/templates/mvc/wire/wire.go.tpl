//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"{{.Module}}/wire/providers"
)

// InitApp 初始化应用，由 wire 自动生成依赖注入代码
// 使用方法：在项目 wire/ 目录下运行 `wire`
func InitApp() (*providers.App, error) {
	wire.Build(
		providers.ProvideConfig,
		providers.ProvideLogger,
		providers.ProvideDatabase,
		providers.ProvideRepository,
		providers.ProvideService,
		providers.ProvideController,
		providers.ProvideRouter,
		wire.Struct(new(providers.App), "*"),
	)
	return &providers.App{}, nil
}
