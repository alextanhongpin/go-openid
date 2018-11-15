package main

import "testing"

func TestEmptyString(t *testing.T) {
	if isEmpty := stringIsEmpty(""); !isEmpty {
		t.Fatalf("want %v, got %v", true, isEmpty)
	}
	if isEmpty := stringIsEmpty(" "); !isEmpty {
		t.Fatalf("want %v, got %v", true, isEmpty)
	}
}
