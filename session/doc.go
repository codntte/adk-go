// Package session provides in-memory session management for ADK agents.
//
// A session represents a single, stateful interaction context between a user
// and an agent. Each session is identified by a unique UUID and can carry
// arbitrary string metadata (e.g. user identity, locale, or conversation
// flags) that persists for the lifetime of the session.
//
// # Basic usage
//
//	store := session.NewStore()
//
//	// Start a new session for an agent.
//	sess := store.Create("my-agent")
//
//	// Attach metadata.
//	_ = store.SetMeta(sess.ID, "user", "alice")
//
//	// Retrieve the session later.
//	s, err := store.Get(sess.ID)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Clean up when done.
//	_ = store.Delete(s.ID)
//
// The Store is safe for concurrent use by multiple goroutines.
package session
