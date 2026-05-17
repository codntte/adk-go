// Package flow provides a lightweight sequential pipeline for chaining
// discrete processing steps.
//
// Each step is a function that receives the accumulated state map and returns
// a map of new or updated keys to merge into that state. Steps are executed
// in the order they are added. If any step returns an error the pipeline halts
// immediately and a [StepError] is returned, wrapping the original error and
// identifying the failing step by name.
//
// Example:
//
//	f := flow.New().
//		Add("enrich", enrichStep).
//		Add("validate", validateStep).
//		Add("store", storeStep)
//
//	result, err := f.Run(ctx, map[string]any{"userID": "u123"})
//	if err != nil {
//		log.Fatal(err)
//	}
//
Flows are safe for concurrent use: multiple goroutines may call Run
simultaneously. Steps themselves are responsible for their own concurrency
safety.
package flow
