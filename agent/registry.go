package agent

import (
	"errors"
	"sync"
)

// Registry holds a collection of named agents.
type Registry struct {
	mu     sync.RWMutex
	agents map[string]*Agent
}

// NewRegistry creates a new, empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]*Agent),
	}
}

// Register adds an agent to the registry.
// Returns an error if an agent with the same name already exists.
func (r *Registry) Register(a *Agent) error {
	if a == nil {
		return errors.New("registry: agent must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.agents[a.Name()]; exists {
		return errors.New("registry: agent already registered: " + a.Name())
	}
	r.agents[a.Name()] = a
	return nil
}

// Get retrieves an agent by name.
func (r *Registry) Get(name string) (*Agent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.agents[name]
	return a, ok
}

// Remove deletes an agent from the registry by name.
func (r *Registry) Remove(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.agents[name]; !ok {
		return false
	}
	delete(r.agents, name)
	return true
}

// Len returns the number of registered agents.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.agents)
}
