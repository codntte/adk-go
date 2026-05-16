// Package runner ties together the session, event, memory, and artifact
// stores to orchestrate agent conversation turns.
//
// # Overview
//
// A Runner is constructed with a session store, an event store, and a
// Handler function. On each call to Run the runner:
//
//  1. Verifies the session exists.
//  2. Records the user input as an event.
//  3. Calls the Handler to produce a response.
//  4. Records the agent response as an event.
//  5. Returns the response to the caller.
//
// # Example
//
//	sessions := session.NewStore()
//	events   := event.NewStore()
//
//	handler := func(ctx context.Context, input string) (string, error) {
//		return "Hello, " + input, nil
//	}
//
//	r, err := runner.New(sessions, events, handler)
//	if err != nil { /* handle */ }
//
//	sid, _ := sessions.Create()
//	out, err := r.Run(ctx, sid, "world")
//	// out == "Hello, world"
package runner
