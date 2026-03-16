package models

import "testing"

func TestCompetitionLevelFromCompetitors(t *testing.T) {
	tests := []struct {
		name     string
		count    int32
		expected string
	}{
		{name: "low lower bound", count: 0, expected: CompetitionLevelLow},
		{name: "low upper bound", count: 10, expected: CompetitionLevelLow},
		{name: "medium lower bound", count: 11, expected: CompetitionLevelMedium},
		{name: "medium upper bound", count: 30, expected: CompetitionLevelMedium},
		{name: "high lower bound", count: 31, expected: CompetitionLevelHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompetitionLevelFromCompetitors(tt.count)
			if got != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
