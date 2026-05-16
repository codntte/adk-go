// Package runner provides a simple agent runner that ties together
// session, event, memory, and artifact stores to execute agent turns.
package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adk-go/adk-go/event"
	"github.com/adk-go/adk-go/session"
)

// Handler is a function that processes an input and returns an output.
type Handler func(ctx context.Context, input string) (string, error)

// Runner orchestrates a single agent conversation turn.
type Runner struct {
	sessions *session.Store
	events   *event.Store
	handler  Handler
}

// New creates a new Runner with the given stores and handler.
func New(sessions *session.Store, events *event.Store, handler Handler) (*Runner, error) {
	if sessions == nil {
		return nil, errors.New("runner: sessions store must not be nil")
	}
	if events == nil {
		return nil, errors.New("runner: events store must not be nil")
	}
	if handler == nil {
		return nil, errors.New("runner: handler must not be nil")
	}
	return &Runner{
		sessions: sessions,
		events:   events,
		handler:  handler,
	}, nil
}

// Run executes a single turn for the given session, appending input and output
// events and returning the handler's response.
func (r *Runner) Run(ctx context.Context, sessionID, input string) (string, error) {
	if _, err := r.sessions.Get(sessionID); err != nil {
		return "", fmt.Errorf("runner: session %q not found: %w", sessionID, err)
	}

	inputEvent := event.Event{
		ID:        fmt.Sprintf("%s-in-%d", sessionID, time.Now().UnixNano()),
		SessionID: sessionID,
		Role:      "user",
		Content:   input,
		CreatedAt: time.Now(),
	}
	if err := r.events.Append(inputEvent); err != nil {
		return "", fmt.Errorf("runner: failed to append input event: %w", err)
	}

	output, err := r.handler(ctx, input)
	if err != nil {
		return "", fmt.Errorf("runner: handler error: %w", err)
	}

	outputEvent := event.Event{
		ID:        fmt.Sprintf("%s-out-%d", sessionID, time.Now().UnixNano()),
		SessionID: sessionID,
		Role:      "agent",
		Content:   output,
		CreatedAt: time.Now(),
	}
	if err := r.events.Append(outputEvent); err != nil {
		return "", fmt.Errorf("runner: failed to append output event: %w", err)
	}

	return output, nil
}

// History returns all events for the given session.
func (r *Runner) History(sessionID string) ([]event.Event, error) {
	return r.events.List(sessionID)
}
