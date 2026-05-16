package runner_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adk-go/adk-go/event"
	"github.com/adk-go/adk-go/runner"
	"github.com/adk-go/adk-go/session"
)

func echoHandler(_ context.Context, input string) (string, error) {
	return "echo: " + input, nil
}

func errHandler(_ context.Context, _ string) (string, error) {
	return "", errors.New("handler failure")
}

func newRunner(t *testing.T, h runner.Handler) (*runner.Runner, *session.Store, *event.Store) {
	t.Helper()
	sessions := session.NewStore()
	events := event.NewStore()
	r, err := runner.New(sessions, events, h)
	if err != nil {
		t.Fatalf("runner.New: %v", err)
	}
	return r, sessions, events
}

func TestRun_Success(t *testing.T) {
	r, sessions, _ := newRunner(t, echoHandler)
	sid, _ := sessions.Create()

	out, err := r.Run(context.Background(), sid, "hello")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if out != "echo: hello" {
		t.Errorf("expected %q, got %q", "echo: hello", out)
	}
}

func TestRun_AppendsEvents(t *testing.T) {
	r, sessions, events := newRunner(t, echoHandler)
	sid, _ := sessions.Create()

	_, _ = r.Run(context.Background(), sid, "ping")

	evs, err := events.List(sid)
	if err != nil {
		t.Fatalf("events.List: %v", err)
	}
	if len(evs) != 2 {
		t.Fatalf("expected 2 events, got %d", len(evs))
	}
	if evs[0].Role != "user" {
		t.Errorf("expected first event role %q, got %q", "user", evs[0].Role)
	}
	if evs[1].Role != "agent" {
		t.Errorf("expected second event role %q, got %q", "agent", evs[1].Role)
	}
}

func TestRun_UnknownSession(t *testing.T) {
	r, _, _ := newRunner(t, echoHandler)

	_, err := r.Run(context.Background(), "no-such-session", "hi")
	if err == nil {
		t.Fatal("expected error for unknown session, got nil")
	}
}

func TestRun_HandlerError(t *testing.T) {
	r, sessions, _ := newRunner(t, errHandler)
	sid, _ := sessions.Create()

	_, err := r.Run(context.Background(), sid, "boom")
	if err == nil {
		t.Fatal("expected error from handler, got nil")
	}
}

func TestNew_NilArgs(t *testing.T) {
	sessions := session.NewStore()
	events := event.NewStore()

	if _, err := runner.New(nil, events, echoHandler); err == nil {
		t.Error("expected error for nil sessions")
	}
	if _, err := runner.New(sessions, nil, echoHandler); err == nil {
		t.Error("expected error for nil events")
	}
	if _, err := runner.New(sessions, events, nil); err == nil {
		t.Error("expected error for nil handler")
	}
}
