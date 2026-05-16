package event

import (
	"sync"
	"time"
)

// Type represents the kind of event.
type Type string

const (
	TypeMessage  Type = "message"
	TypeToolCall Type = "tool_call"
	TypeToolResult Type = "tool_result"
	TypeError    Type = "error"
)

// Event represents a discrete occurrence within an agent session.
type Event struct {
	ID        string
	SessionID string
	Type      Type
	Payload   map[string]any
	CreatedAt time.Time
}

// Store holds events indexed by session ID.
type Store struct {
	mu     sync.RWMutex
	events map[string][]*Event
}

// NewStore creates and returns a new event Store.
func NewStore() *Store {
	return &Store{
		events: make(map[string][]*Event),
	}
}

// Append adds an event to the store under its session ID.
func (s *Store) Append(e *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.SessionID] = append(s.events[e.SessionID], e)
}

// ListBySession returns all events for the given session ID.
func (s *Store) ListBySession(sessionID string) []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	evts := s.events[sessionID]
	if evts == nil {
		return []*Event{}
	}
	result := make([]*Event, len(evts))
	copy(result, evts)
	return result
}

// DeleteBySession removes all events associated with a session ID.
func (s *Store) DeleteBySession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, sessionID)
}

// Len returns the total number of events stored across all sessions.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total := 0
	for _, evts := range s.events {
		total += len(evts)
	}
	return total
}
