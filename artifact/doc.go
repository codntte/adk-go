// Package artifact provides an in-memory store for binary and text blobs
// (artifacts) that are associated with agent sessions.
//
// An [Artifact] is a named, typed blob — for example an uploaded file, a
// generated image, or a structured JSON document — that an agent may produce
// or consume during a session.  Each artifact belongs to exactly one session
// and is identified by a unique ID.
//
// # Store
//
// [Store] is the primary type in this package.  It is safe for concurrent use
// by multiple goroutines.
//
//	s := artifact.NewStore()
//
//	err := s.Set(&artifact.Artifact{
//		ID:        "art-001",
//		SessionID: "sess-abc",
//		Name:      "report.pdf",
//		MIMEType:  "application/pdf",
//		Data:      pdfBytes,
//	})
//
//	a, err := s.Get("art-001")
//
// When a session ends, all of its artifacts can be removed in one call:
//
//	n := s.DeleteBySession("sess-abc")
//
// # Errors
//
// [ErrNotFound] is returned by [Store.Get] and [Store.Delete] when the
// requested artifact does not exist.
package artifact
