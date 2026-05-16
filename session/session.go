// Package session provides session management for ADK agents,
// allowing state to be tracked across multiple interactions.
package session

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ErrSessionNotFound is returned when a session cannot be located.
var ErrSessionNotFound = errors.New("session: not found")

// Session represents a single agent interaction session.
type Session struct {
	ID        string
	AgentID   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Metadata  map[string]string
}

// Store manages the lifecycle of sessions in memory.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewStore creates and returns a new session Store.
func NewStore() *Store {
	return &Store{
		sessions: make(map[string]*Session),
	}
}

// Create initialises a new session for the given agentID and returns it.
func (s *Store) Create(agentID string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	sess := &Session{
		ID:        uuid.NewString(),
		AgentID:   agentID,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  make(map[string]string),
	}
	s.sessions[sess.ID] = sess
	return sess
}

// Get retrieves a session by its ID.
func (s *Store) Get(id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return sess, nil
}

// SetMeta sets a metadata key/value pair on an existing session.
func (s *Store) SetMeta(id, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[id]
	if !ok {
		return ErrSessionNotFound
	}
	sess.Metadata[key] = value
	sess.UpdatedAt = time.Now().UTC()
	return nil
}

// Delete removes a session from the store.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[id]; !ok {
		return ErrSessionNotFound
	}
	delete(s.sessions, id)
	return nil
}

// Len returns the number of active sessions.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}
