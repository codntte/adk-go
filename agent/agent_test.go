package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adk-go/adk-go/agent"
)

func echoHandler(_ context.Context, input string) (string, error) {
	return input, nil
}

func errHandler(_ context.Context, _ string) (string, error) {
	return "", errors.New("handler error")
}

func TestNew_Success(t *testing.T) {
	a, err := agent.New("test-agent", "a test agent", echoHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "test-agent" {
		t.Errorf("expected name %q, got %q", "test-agent", a.Name())
	}
	if a.Description() != "a test agent" {
		t.Errorf("expected description %q, got %q", "a test agent", a.Description())
	}
}

func TestNew_EmptyName(t *testing.T) {
	_, err := agent.New("", "desc", echoHandler)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNew_NilHandler(t *testing.T) {
	_, err := agent.New("agent", "desc", nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestRun_Success(t *testing.T) {
	a, _ := agent.New("echo", "", echoHandler)
	out, err := a.Run(context.Background(), "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello" {
		t.Errorf("expected %q, got %q", "hello", out)
	}
}

func TestRun_HandlerError(t *testing.T) {
	a, _ := agent.New("err-agent", "", errHandler)
	_, err := a.Run(context.Background(), "input")
	if err == nil {
		t.Fatal("expected error from handler")
	}
}

func TestRun_NilContext(t *testing.T) {
	a, _ := agent.New("agent", "", echoHandler)
	_, err := a.Run(nil, "input") //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil context")
	}
}

func TestMeta(t *testing.T) {
	a, _ := agent.New("agent", "", echoHandler)
	a.SetMeta("model", "gemini-pro")
	v, ok := a.GetMeta("model")
	if !ok {
		t.Fatal("expected meta key to exist")
	}
	if v != "gemini-pro" {
		t.Errorf("expected %q, got %q", "gemini-pro", v)
	}
	_, ok = a.GetMeta("missing")
	if ok {
		t.Fatal("expected missing key to not exist")
	}
}
