//go:build common_integration

package commonadapter

import (
    "context"
    "fmt"

    _ "github.com/yourname/common"
    "github.com/yourname/go-start/pkg/spec"
)

// CommonAdapter 为启用 common 集成时的适配器实现
// 该文件受构建标签控制，仅在启用 `-tags common_integration` 时参与构建。
type CommonAdapter struct{}

// Init 初始化 common 集成环境
// 具体初始化逻辑依赖 common 包的实现，当前为占位返回。
func (a *CommonAdapter) Init(ctx context.Context, cfg Config) error {
    _ = ctx
    _ = cfg
    return nil
}

// Generate 调用 common 包完成生成流程
// 具体调用需在 common 包就绪后补充，此处占位返回明确错误。
func (a *CommonAdapter) Generate(ctx context.Context, s *spec.Spec, outputDir string) (Result, error) {
    _ = ctx
    _ = s
    _ = outputDir
    return Result{Success: false, Message: "common 集成未完成"}, fmt.Errorf("common 集成占位: %w", ErrCommonUnavailable)
}

// Validate 使用 common 或本地规则进行校验
// 当前占位实现返回 nil，实际逻辑待 common 就绪后填充。
func (a *CommonAdapter) Validate(ctx context.Context, s *spec.Spec) error {
    _ = ctx
    _ = s
    return nil
}