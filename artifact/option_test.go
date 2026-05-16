package artifact_test

import (
	"testing"

	"github.com/adk-go/adk-go/artifact"
)

func TestWithMaxPerSession_UnderLimit(t *testing.T) {
	s := artifact.NewStoreWithOptions(artifact.WithMaxPerSession(3))
	for i, id := range []string{"a1", "a2", "a3"} {
		_ = i
		if err := s.Set(newArtifact(id, "s1", id+".txt")); err != nil {
			t.Fatalf("Set %s: %v", id, err)
		}
	}
	if s.Len() != 3 {
		t.Errorf("Len = %d, want 3", s.Len())
	}
}

func TestWithMaxPerSession_AtLimit(t *testing.T) {
	s := artifact.NewStoreWithOptions(artifact.WithMaxPerSession(2))
	_ = s.Set(newArtifact("a1", "s1", "f1.txt"))
	_ = s.Set(newArtifact("a2", "s1", "f2.txt"))
	err := s.Set(newArtifact("a3", "s1", "f3.txt"))
	if err != artifact.ErrLimitExceeded {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
	if s.Len() != 2 {
		t.Errorf("Len = %d, want 2", s.Len())
	}
}

func TestWithMaxPerSession_DifferentSessions(t *testing.T) {
	s := artifact.NewStoreWithOptions(artifact.WithMaxPerSession(1))
	if err := s.Set(newArtifact("a1", "s1", "f1.txt")); err != nil {
		t.Fatalf("Set s1/a1: %v", err)
	}
	// Different session — should succeed even though s1 is at its limit.
	if err := s.Set(newArtifact("a2", "s2", "f2.txt")); err != nil {
		t.Fatalf("Set s2/a2: %v", err)
	}
	if s.Len() != 2 {
		t.Errorf("Len = %d, want 2", s.Len())
	}
}

func TestWithMaxPerSession_UpdateDoesNotCount(t *testing.T) {
	s := artifact.NewStoreWithOptions(artifact.WithMaxPerSession(1))
	a := newArtifact("a1", "s1", "f1.txt")
	_ = s.Set(a)
	// Re-setting the same ID should not be treated as a new artifact.
	a.Name = "updated.txt"
	if err := s.Set(a); err != nil {
		t.Fatalf("update Set: %v", err)
	}
}
