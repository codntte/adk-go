// Package event provides an in-memory store for tracking discrete events
// that occur within agent sessions in adk-go.
//
// Events capture meaningful occurrences such as messages exchanged between
// the user and the agent, tool calls initiated by the agent, tool results
// returned to the agent, and errors encountered during processing.
//
// # Usage
//
// Create a new store, then append events as they occur:
//
//	store := event.NewStore()
//
//	store.Append(&event.Event{
//		ID:        "evt-001",
//		SessionID: "sess-abc",
//		Type:      event.TypeMessage,
//		Payload:   map[string]any{"text": "Hello, agent!"},
//		CreatedAt: time.Now(),
//	})
//
//	events := store.ListBySession("sess-abc")
//
// Events are stored per session and can be retrieved or deleted by session ID.
package event
