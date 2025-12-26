package commonadapter

import (
	"context"
	"time"

	"github.com/yourname/go-start/pkg/spec"
	"go.uber.org/zap"
)

// Config 表示适配器的基础配置
// 包含日志、超时等可选设置，用于初始化工·具包调用环境。
type Config struct {
	Logger  *zap.Logger
	Timeout time.Duration
}

// Result 表示一次生成流程的结果
// 包含成功标记与消息描述，便于上层展示与调试。
type Result struct {
	Success bool
	Message string
}

// CommonTools 定义 go-start 使用 common 工具包的能力边界
// 通过接口隔离具体实现，便于测试与未来替换（如插件化）。
type CommonTools interface {
	// Init 初始化工具包调用环境
	// 入参包含上下文与配置，返回错误以指示初始化失败。
	Init(ctx context.Context, cfg Config) error

	// Generate 根据规范生成工程产物
	// 入参为解析后的规范与输出目录，返回结果与错误。
	Generate(ctx context.Context, s *spec.Spec, outputDir string) (Result, error)

	// Validate 对规范进行基础校验
	// 返回错误以指示规范不合法或环境不满足。
	Validate(ctx context.Context, s *spec.Spec) error
}

// NewAdapter 创建默认适配器实例
// 默认使用本地无依赖的 NoopAdapter，实现与 common 解耦的最小能力。
func NewAdapter() CommonTools {
	return &NoopAdapter{}
}

// WithLogger 为适配器配置日志
// 若未提供日志，则在实现内部可能使用空日志以避免空指针。
func WithLogger(a CommonTools, logger *zap.Logger) CommonTools {
	// 该帮助函数保留用于未来需要在构造期注入日志的实现。
	_ = logger
	return a
}

// Abilities 返回默认的各类能力接口实现
// 在未启用 common 集成时返回 Noop 实现，高级工程师可通过构建标签替换。
func Abilities() (Auth, Cache, Audit, Idempotency, Lock, EventBus) {
	return &NoopAuth{}, &NoopCache{}, &NoopAudit{}, &NoopIdempotency{}, &NoopLock{}, &NoopEventBus{}
}
