package main

import "testing"

func TestGenerateIdeas(t *testing.T) {
	ideas := generateIdeas("AI meeting assistant")
	if len(ideas) != len(niches) {
		t.Fatalf("expected %d ideas, got %d", len(niches), len(ideas))
	}
	if ideas[0] != "AI meeting assistant for recruiters" {
		t.Fatalf("unexpected first idea: %s", ideas[0])
	}
}
