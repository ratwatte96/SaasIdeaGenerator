package models

const (
	CompetitionLevelLow    = "low"
	CompetitionLevelMedium = "medium"
	CompetitionLevelHigh   = "high"
)

// CompetitionLevelFromCompetitors maps competitor_count to deterministic buckets.
// low: 0-10, medium: 11-30, high: 31+
func CompetitionLevelFromCompetitors(competitorCount int32) string {
	if competitorCount <= 10 {
		return CompetitionLevelLow
	}
	if competitorCount <= 30 {
		return CompetitionLevelMedium
	}
	return CompetitionLevelHigh
}
