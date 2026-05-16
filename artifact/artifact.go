package artifact

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when an artifact is not found in the store.
var ErrNotFound = errors.New("artifact: not found")

// Artifact represents a named binary or text blob associated with a session.
type Artifact struct {
	ID        string
	SessionID string
	Name      string
	MIMEType  string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Store is a thread-safe in-memory store for artifacts.
type Store struct {
	mu    sync.RWMutex
	items map[string]*Artifact // keyed by Artifact.ID
}

// NewStore returns an initialised, empty Store.
func NewStore() *Store {
	return &Store{
		items: make(map[string]*Artifact),
	}
}

// Set inserts or replaces an artifact. The artifact must have a non-empty ID
// and SessionID.
func (s *Store) Set(a *Artifact) error {
	if a.ID == "" {
		return errors.New("artifact: ID must not be empty")
	}
	if a.SessionID == "" {
		return errors.New("artifact: SessionID must not be empty")
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.items[a.ID]; ok {
		a.CreatedAt = existing.CreatedAt
	} else {
		a.CreatedAt = now
	}
	a.UpdatedAt = now
	copy := *a
	s.items[a.ID] = &copy
	return nil
}

// Get returns the artifact with the given ID or ErrNotFound.
func (s *Store) Get(id string) (*Artifact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *a
	return &copy, nil
}

// ListBySession returns all artifacts that belong to the given session.
func (s *Store) ListBySession(sessionID string) []*Artifact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Artifact
	for _, a := range s.items {
		if a.SessionID == sessionID {
			copy := *a
			out = append(out, &copy)
		}
	}
	return out
}

// Delete removes the artifact with the given ID. It returns ErrNotFound if
// the artifact does not exist.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}

// DeleteBySession removes all artifacts belonging to sessionID and returns the
// number of entries removed.
func (s *Store) DeleteBySession(sessionID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for id, a := range s.items {
		if a.SessionID == sessionID {
			delete(s.items, id)
			count++
		}
	}
	return count
}

// Len returns the total number of artifacts in the store.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
