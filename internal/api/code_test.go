package main

import (
	"testing"
	"time"
)

func TestCodeFactory(t *testing.T) {
	code := Code{
		ID:        "1",
		TTL:       1 * time.Second,
		CreatedAt: time.Unix(0, 0),
	}
	cf := NewCodeBuilder(code)
	result := cf.Build()
	if id := result.ID; id != code.ID {
		t.Fatalf("want %v, got %v", code.ID, id)
	}
	if ttl := result.TTL; ttl != code.TTL {
		t.Fatalf("want %v, got %v", code.TTL, ttl)
	}
	if createdAt := result.CreatedAt; createdAt != code.CreatedAt {
		t.Fatalf("want %v, got %v", code.CreatedAt, createdAt)
	}
}
