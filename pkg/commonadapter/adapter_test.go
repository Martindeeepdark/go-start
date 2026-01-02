package commonadapter

import (
	"context"
	"testing"

	"github.com/Martindeeepdark/go-start/pkg/spec"
)

// TestNoopAdapterValidate 验证默认实现的基础校验逻辑
func TestNoopAdapterValidate(t *testing.T) {
	a := &NoopAdapter{}
	ctx := context.Background()

	s := &spec.Spec{
		Name:    "ExampleAPI",
		Project: spec.ProjectConfig{Module: "github.com/Martindeeepdark/example"},
	}

	if err := a.Validate(ctx, s); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

// TestNoopAdapterGenerate 验证默认实现在未集成 common 时的错误返回
func TestNoopAdapterGenerate(t *testing.T) {
	a := &NoopAdapter{}
	ctx := context.Background()
	s := &spec.Spec{
		Name:    "ExampleAPI",
		Project: spec.ProjectConfig{Module: "github.com/Martindeeepdark/example"},
	}

	_, err := a.Generate(ctx, s, "/tmp/output")
	if err == nil {
		t.Fatalf("Generate() expected error, got nil")
	}
	if err != ErrCommonUnavailable {
		t.Fatalf("Generate() expected ErrCommonUnavailable, got %v", err)
	}
}
