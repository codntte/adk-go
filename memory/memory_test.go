package memory_test

import (
	"testing"

	"github.com/adk-go/adk-go/memory"
)

func TestStore_SetAndGet(t *testing.T) {
	s := memory.NewStore()

	s.Set("foo", "bar")

	entry, err := s.Get("foo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if entry.Value != "bar" {
		t.Errorf("expected value 'bar', got %v", entry.Value)
	}
	if entry.Key != "foo" {
		t.Errorf("expected key 'foo', got %v", entry.Key)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	s := memory.NewStore()

	_, err := s.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != memory.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_Update(t *testing.T) {
	s := memory.NewStore()

	s.Set("key", "initial")
	s.Set("key", "updated")

	entry, err := s.Get("key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Value != "updated" {
		t.Errorf("expected 'updated', got %v", entry.Value)
	}
	if !entry.UpdatedAt.After(entry.CreatedAt) && entry.UpdatedAt != entry.CreatedAt {
		t.Error("expected UpdatedAt to be set on update")
	}
}

func TestStore_Delete(t *testing.T) {
	s := memory.NewStore()

	s.Set("todelete", 42)
	s.Delete("todelete")

	_, err := s.Get("todelete")
	if err != memory.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestStore_Len(t *testing.T) {
	s := memory.NewStore()

	if s.Len() != 0 {
		t.Errorf("expected 0, got %d", s.Len())
	}

	s.Set("a", 1)
	s.Set("b", 2)

	if s.Len() != 2 {
		t.Errorf("expected 2, got %d", s.Len())
	}
}

func TestStore_Keys(t *testing.T) {
	s := memory.NewStore()

	s.Set("x", 1)
	s.Set("y", 2)

	keys := s.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
