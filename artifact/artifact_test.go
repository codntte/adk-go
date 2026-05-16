package artifact_test

import (
	"testing"

	"github.com/adk-go/adk-go/artifact"
)

func newArtifact(id, sessionID, name string) *artifact.Artifact {
	return &artifact.Artifact{
		ID:        id,
		SessionID: sessionID,
		Name:      name,
		MIMEType:  "text/plain",
		Data:      []byte("hello"),
	}
}

func TestStore_SetAndGet(t *testing.T) {
	s := artifact.NewStore()
	a := newArtifact("a1", "s1", "file.txt")
	if err := s.Set(a); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := s.Get("a1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "file.txt" {
		t.Errorf("Name = %q, want %q", got.Name, "file.txt")
	}
}

func TestStore_GetNotFound(t *testing.T) {
	s := artifact.NewStore()
	_, err := s.Get("missing")
	if err != artifact.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_SetEmptyID(t *testing.T) {
	s := artifact.NewStore()
	a := newArtifact("", "s1", "file.txt")
	if err := s.Set(a); err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestStore_Update(t *testing.T) {
	s := artifact.NewStore()
	a := newArtifact("a1", "s1", "file.txt")
	_ = s.Set(a)
	a.Name = "updated.txt"
	_ = s.Set(a)
	got, _ := s.Get("a1")
	if got.Name != "updated.txt" {
		t.Errorf("Name = %q, want %q", got.Name, "updated.txt")
	}
	if got.CreatedAt != got.UpdatedAt && got.CreatedAt.IsZero() {
		t.Error("CreatedAt should be preserved on update")
	}
}

func TestStore_ListBySession(t *testing.T) {
	s := artifact.NewStore()
	_ = s.Set(newArtifact("a1", "s1", "f1.txt"))
	_ = s.Set(newArtifact("a2", "s1", "f2.txt"))
	_ = s.Set(newArtifact("a3", "s2", "f3.txt"))
	list := s.ListBySession("s1")
	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
}

func TestStore_Delete(t *testing.T) {
	s := artifact.NewStore()
	_ = s.Set(newArtifact("a1", "s1", "f1.txt"))
	if err := s.Delete("a1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("Len = %d, want 0", s.Len())
	}
}

func TestStore_DeleteNotFound(t *testing.T) {
	s := artifact.NewStore()
	if err := s.Delete("nope"); err != artifact.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_DeleteBySession(t *testing.T) {
	s := artifact.NewStore()
	_ = s.Set(newArtifact("a1", "s1", "f1.txt"))
	_ = s.Set(newArtifact("a2", "s1", "f2.txt"))
	_ = s.Set(newArtifact("a3", "s2", "f3.txt"))
	n := s.DeleteBySession("s1")
	if n != 2 {
		t.Errorf("deleted %d, want 2", n)
	}
	if s.Len() != 1 {
		t.Errorf("Len = %d, want 1", s.Len())
	}
}

func TestStore_Len(t *testing.T) {
	s := artifact.NewStore()
	if s.Len() != 0 {
		t.Errorf("empty store Len = %d, want 0", s.Len())
	}
	_ = s.Set(newArtifact("a1", "s1", "f1.txt"))
	if s.Len() != 1 {
		t.Errorf("Len = %d, want 1", s.Len())
	}
}
