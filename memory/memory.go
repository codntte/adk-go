// Package memory provides in-memory storage for agent session data.
package memory

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when a session is not found.
var ErrNotFound = errors.New("session not found")

// Entry represents a stored memory entry.
type Entry struct {
	Key       string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Store is a thread-safe in-memory key-value store for agent sessions.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// NewStore creates and returns a new Store instance.
func NewStore() *Store {
	return &Store{
		entries: make(map[string]*Entry),
	}
}

// Set stores a value under the given key.
func (s *Store) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if existing, ok := s.entries[key]; ok {
		existing.Value = value
		existing.UpdatedAt = now
		return
	}
	s.entries[key] = &Entry{
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Get retrieves a value by key. Returns ErrNotFound if the key does not exist.
func (s *Store) Get(key string) (*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.entries[key]
	if !ok {
		return nil, ErrNotFound
	}
	return entry, nil
}

// Delete removes an entry by key.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Keys returns all stored keys.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of stored entries.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
