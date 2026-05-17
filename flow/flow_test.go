package flow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adk-go/adk-go/flow"
)

func addStep(key string, value any) flow.Step {
	return func(_ context.Context, input map[string]any) (map[string]any, error) {
		return map[string]any{key: value}, nil
	}
}

func TestFlow_RunSuccess(t *testing.T) {
	f := flow.New().
		Add("step1", addStep("a", 1)).
		Add("step2", addStep("b", 2))

	out, err := f.Run(context.Background(), map[string]any{"init": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["a"] != 1 || out["b"] != 2 || out["init"] != true {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestFlow_RunStepError(t *testing.T) {
	sentinel := errors.New("boom")
	f := flow.New().
		Add("step1", addStep("a", 1)).
		Add("fail", func(_ context.Context, _ map[string]any) (map[string]any, error) {
			return nil, sentinel
		})

	_, err := f.Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var stepErr *flow.StepError
	if !errors.As(err, &stepErr) {
		t.Fatalf("expected StepError, got %T", err)
	}
	if stepErr.Step != "fail" {
		t.Errorf("expected step name 'fail', got %q", stepErr.Step)
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error to be unwrapped")
	}
}

func TestFlow_RunCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	f := flow.New().Add("step1", addStep("a", 1))
	_, err := f.Run(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestFlow_RunNilContext(t *testing.T) {
	f := flow.New().Add("step1", addStep("a", 1))
	//nolint:staticcheck
	_, err := f.Run(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil context")
	}
}

func TestFlow_Len(t *testing.T) {
	f := flow.New()
	if f.Len() != 0 {
		t.Errorf("expected 0, got %d", f.Len())
	}
	f.Add("s1", addStep("x", 1)).Add("s2", addStep("y", 2))
	if f.Len() != 2 {
		t.Errorf("expected 2, got %d", f.Len())
	}
}

func TestFlow_PanicOnEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty step name")
		}
	}()
	flow.New().Add("", addStep("x", 1))
}
