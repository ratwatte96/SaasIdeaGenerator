package store

import "time"

type Product struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Category         string    `json:"category"`
	Source           string    `json:"source"`
	SourceExternalID string    `json:"source_external_id"`
	Upvotes          int       `json:"upvotes"`
	CompetitorCount  int       `json:"competitor_count"`
	CreatedAt        time.Time `json:"created_at"`
}

type Idea struct {
	ID               string    `json:"id"`
	IdeaText         string    `json:"idea_text"`
	SourceProductID  string    `json:"source_product_id"`
	DemandScore      float64   `json:"demand_score"`
	CompetitionLevel string    `json:"competition_level"`
	CreatedAt        time.Time `json:"created_at"`
	ProductName      string    `json:"product_name"`
	ProductCategory  string    `json:"product_category"`
}
