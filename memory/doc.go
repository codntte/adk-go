// Package memory implements a lightweight, thread-safe in-memory store
// for managing agent session state within the adk-go framework.
//
// # Overview
//
// The Store type provides basic CRUD operations for key-value pairs used
// by agents during their lifecycle. Each entry tracks creation and update
// timestamps, making it suitable for session management and short-lived
// context storage.
//
// # Usage
//
//	store := memory.NewStore()
//
//	// Store a value
//	store.Set("session_id", "abc-123")
//
//	// Retrieve a value
//	entry, err := store.Get("session_id")
//	if err != nil {
//		// handle memory.ErrNotFound
//	}
//
//	// Remove a value
//	store.Delete("session_id")
//
// # Concurrency
//
// All operations on Store are safe for concurrent use by multiple goroutines.
package memory
