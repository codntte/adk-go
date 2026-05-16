package session_test

import (
	"testing"

	"github.com/google/adk-go/session"
)

func TestStore_Create(t *testing.T) {
	store := session.NewStore()
	sess := store.Create("agent-1")

	if sess.ID == "" {
		t.Fatal("expected non-empty session ID")
	}
	if sess.AgentID != "agent-1" {
		t.Fatalf("expected agentID %q, got %q", "agent-1", sess.AgentID)
	}
	if store.Len() != 1 {
		t.Fatalf("expected Len 1, got %d", store.Len())
	}
}

func TestStore_Get(t *testing.T) {
	store := session.NewStore()
	sess := store.Create("agent-2")

	got, err := store.Get(sess.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != sess.ID {
		t.Fatalf("expected ID %q, got %q", sess.ID, got.ID)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	store := session.NewStore()

	_, err := store.Get("nonexistent")
	if err != session.ErrSessionNotFound {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestStore_SetMeta(t *testing.T) {
	store := session.NewStore()
	sess := store.Create("agent-3")

	if err := store.SetMeta(sess.ID, "user", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := store.Get(sess.ID)
	if got.Metadata["user"] != "alice" {
		t.Fatalf("expected metadata user=alice, got %q", got.Metadata["user"])
	}
}

func TestStore_SetMetaNotFound(t *testing.T) {
	store := session.NewStore()

	err := store.SetMeta("missing", "k", "v")
	if err != session.ErrSessionNotFound {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestStore_Delete(t *testing.T) {
	store := session.NewStore()
	sess := store.Create("agent-4")

	if err := store.Delete(sess.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Len() != 0 {
		t.Fatalf("expected Len 0 after delete, got %d", store.Len())
	}
}

func TestStore_DeleteNotFound(t *testing.T) {
	store := session.NewStore()

	err := store.Delete("ghost")
	if err != session.ErrSessionNotFound {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}
