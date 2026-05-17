// Package agent provides primitives for defining and managing AI agents
// within the adk-go framework.
//
// An Agent encapsulates a name, a human-readable description, and a Handler
// function that processes a string input and returns a string output. Agents
// are safe for concurrent use.
//
// Basic usage:
//
//	import "github.com/adk-go/adk-go/agent"
//
//	a, err := agent.New("my-agent", "does something useful", func(ctx context.Context, input string) (string, error) {
//	    return "response: " + input, nil
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	output, err := a.Run(ctx, "hello")
//
// A Registry can be used to manage multiple agents by name:
//
//	reg := agent.NewRegistry()
//	reg.Register(a)
//	found, ok := reg.Get("my-agent")
package agent
