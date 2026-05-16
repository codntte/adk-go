package event_test

import (
	"testing"
	"time"

	"github.com/adk-go/adk-go/event"
)

func newEvent(id, sessionID string, t event.Type) *event.Event {
	return &event.Event{
		ID:        id,
		SessionID: sessionID,
		Type:      t,
		Payload:   map[string]any{"text": "hello"},
		CreatedAt: time.Now(),
	}
}

func TestStore_AppendAndList(t *testing.T) {
	s := event.NewStore()
	s.Append(newEvent("e1", "sess-1", event.TypeMessage))
	s.Append(newEvent("e2", "sess-1", event.TypeToolCall))

	evts := s.ListBySession("sess-1")
	if len(evts) != 2 {
		t.Fatalf("expected 2 events, got %d", len(evts))
	}
}

func TestStore_ListEmpty(t *testing.T) {
	s := event.NewStore()
	evts := s.ListBySession("nonexistent")
	if len(evts) != 0 {
		t.Fatalf("expected 0 events, got %d", len(evts))
	}
}

func TestStore_DeleteBySession(t *testing.T) {
	s := event.NewStore()
	s.Append(newEvent("e1", "sess-1", event.TypeMessage))
	s.Append(newEvent("e2", "sess-2", event.TypeError))

	s.DeleteBySession("sess-1")

	if len(s.ListBySession("sess-1")) != 0 {
		t.Fatal("expected sess-1 events to be deleted")
	}
	if len(s.ListBySession("sess-2")) != 1 {
		t.Fatal("expected sess-2 events to remain")
	}
}

func TestStore_Len(t *testing.T) {
	s := event.NewStore()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Append(newEvent("e1", "sess-1", event.TypeMessage))
	s.Append(newEvent("e2", "sess-1", event.TypeToolResult))
	s.Append(newEvent("e3", "sess-2", event.TypeMessage))
	if s.Len() != 3 {
		t.Fatalf("expected 3, got %d", s.Len())
	}
}

func TestStore_IsolatedSessions(t *testing.T) {
	s := event.NewStore()
	s.Append(newEvent("e1", "sess-A", event.TypeMessage))
	s.Append(newEvent("e2", "sess-B", event.TypeMessage))

	if len(s.ListBySession("sess-A")) != 1 {
		t.Fatal("sess-A should have exactly 1 event")
	}
	if len(s.ListBySession("sess-B")) != 1 {
		t.Fatal("sess-B should have exactly 1 event")
	}
}
