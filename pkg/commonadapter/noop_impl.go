package commonadapter

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/Martindeeepdark/go-start/pkg/spec"
)

// NoopAdapter 为默认的本地实现
// 不依赖 external common 包，提供基本的校验与占位行为。
type NoopAdapter struct{}

// Init 初始化本地适配器
// 该实现不做外部连接，仅模拟耗时以保持行为一致性。
func (a *NoopAdapter) Init(ctx context.Context, cfg Config) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(10 * time.Millisecond):
        return nil
    }
}

// Generate 执行占位生成逻辑
// 在未启用 common 集成时返回明确错误，以提示使用者。
func (a *NoopAdapter) Generate(ctx context.Context, s *spec.Spec, outputDir string) (Result, error) {
    if s == nil {
        return Result{}, fmt.Errorf("spec 不能为空: %w", errors.New("nil spec"))
    }
    if err := a.Validate(ctx, s); err != nil {
        return Result{}, err
    }
    return Result{Success: false, Message: "common 未集成，生成功能不可用"}, ErrCommonUnavailable
}

// Validate 对规范进行轻量校验
// 校验关键字段，避免进入后续流程产生不可预期错误。
func (a *NoopAdapter) Validate(ctx context.Context, s *spec.Spec) error {
    if s == nil {
        return fmt.Errorf("规范为空")
    }
    if s.Name == "" {
        return fmt.Errorf("缺少 API 名称")
    }
    if s.Project.Module == "" {
        return fmt.Errorf("缺少项目模块名")
    }
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        return nil
    }
}