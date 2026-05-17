package agent

import (
	"context"
	"errors"
	"sync"
)

// Handler is a function that processes a message and returns a response.
type Handler func(ctx context.Context, input string) (string, error)

// Agent represents an AI agent with a name, description, and handler.
type Agent struct {
	mu          sync.RWMutex
	name        string
	description string
	handler     Handler
	meta        map[string]string
}

// New creates a new Agent with the given name, description, and handler.
func New(name, description string, handler Handler) (*Agent, error) {
	if name == "" {
		return nil, errors.New("agent: name must not be empty")
	}
	if handler == nil {
		return nil, errors.New("agent: handler must not be nil")
	}
	return &Agent{
		name:        name,
		description: description,
		handler:     handler,
		meta:        make(map[string]string),
	}, nil
}

// Name returns the agent's name.
func (a *Agent) Name() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.name
}

// Description returns the agent's description.
func (a *Agent) Description() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.description
}

// SetMeta sets a metadata key-value pair on the agent.
func (a *Agent) SetMeta(key, value string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.meta[key] = value
}

// GetMeta retrieves a metadata value by key.
func (a *Agent) GetMeta(key string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	v, ok := a.meta[key]
	return v, ok
}

// Run executes the agent's handler with the given context and input.
func (a *Agent) Run(ctx context.Context, input string) (string, error) {
	if ctx == nil {
		return "", errors.New("agent: context must not be nil")
	}
	a.mu.RLock()
	h := a.handler
	a.mu.RUnlock()
	return h(ctx, input)
}
