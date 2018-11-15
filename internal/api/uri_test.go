package main

import "testing"

func TestURIContains(t *testing.T) {
	uris := []string{"a", "b", "c"}
	if yes := URIs(uris).Contains("a"); !yes {
		t.Fatalf("want %v, got %v", true, yes)
	}
	if no := URIs(uris).Contains("z"); no {
		t.Fatalf("want %v, got %v", false, no)
	}
}
