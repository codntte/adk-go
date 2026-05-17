package flow

import (
	"context"
	"errors"
	"sync"
)

// Step represents a single unit of work in a flow.
type Step func(ctx context.Context, input map[string]any) (map[string]any, error)

// Flow executes a sequence of steps, passing output of each step as input to the next.
type Flow struct {
	mu    sync.RWMutex
	steps []namedStep
}

type namedStep struct {
	name string
	fn   Step
}

// New creates a new empty Flow.
func New() *Flow {
	return &Flow{}
}

// Add appends a named step to the flow.
func (f *Flow) Add(name string, step Step) *Flow {
	if name == "" {
		panic("flow: step name must not be empty")
	}
	if step == nil {
		panic("flow: step func must not be nil")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.steps = append(f.steps, namedStep{name: name, fn: step})
	return f
}

// Len returns the number of steps in the flow.
func (f *Flow) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.steps)
}

// Run executes all steps in order. Each step receives the merged output of all
// previous steps plus the original input. Returns the final accumulated state.
func (f *Flow) Run(ctx context.Context, input map[string]any) (map[string]any, error) {
	if ctx == nil {
		return nil, errors.New("flow: context must not be nil")
	}

	f.mu.RLock()
	steps := make([]namedStep, len(f.steps))
	copy(steps, f.steps)
	f.mu.RUnlock()

	state := make(map[string]any, len(input))
	for k, v := range input {
		state[k] = v
	}

	for _, s := range steps {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		out, err := s.fn(ctx, state)
		if err != nil {
			return nil, &StepError{Step: s.name, Err: err}
		}
		for k, v := range out {
			state[k] = v
		}
	}
	return state, nil
}

// StepError wraps an error returned by a named step.
type StepError struct {
	Step string
	Err  error
}

func (e *StepError) Error() string {
	return "flow: step \"" + e.Step + "\": " + e.Err.Error()
}

func (e *StepError) Unwrap() error { return e.Err }
