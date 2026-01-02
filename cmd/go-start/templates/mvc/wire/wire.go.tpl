//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"{{.Module}}/cmd/server/providers"
)

// ============================================
// Wire 注入器
// ============================================
//
// 说明：
// 1. Wire 会扫描这个文件中的函数
// 2. 自动生成 wire_gen.go 文件
// 3. 运行 `wire` 命令生成代码
//
// 使用方法：
//   cd cmd/server/wire
//   wire

// InitApp 初始化应用
//
// Wire 会根据返回类型和 ProviderSet
// 自动构建依赖注入链
func InitApp() (*App, error) {
	wire.Build(
		// 所有 Provider
		providers.ProviderSet,

		// App 结构体
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

// ============================================
// Provider Set：所有 Provider 的集合
// ============================================

// ProviderSet 包含所有依赖提供者
//
// 说明：
// - 将所有 Provider 组织在一起
// - 便于管理和复用
// - Wire 会根据类型自动匹配
var ProviderSet = wire.NewSet(
	providers.ProvideConfig,
	providers.ProvideLogger,
	providers.ProvideDatabase,
	providers.ProvideCache,
	providers.ProvideRepository,
	providers.ProvideService,
	providers.ProvideController,
	providers.ProvideRouter,
)
