package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"saasideagenerator/backend/internal/config"
	"saasideagenerator/backend/internal/db"
	"saasideagenerator/backend/internal/store"
	"saasideagenerator/backend/models"
)

var niches = []string{"recruiters", "sales teams", "lawyers", "consultants", "real estate agents", "content creators"}

type sourceProduct struct {
	Name            string
	Tagline         string
	Category        string
	Upvotes         int
	ExternalID      string
	CompetitorCount int
}

func main() {
	cfg := config.Load()
	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer database.Close()

	s := store.New(database)
	ctx := context.Background()
	products := mockProducts()
	if os.Getenv("COLLECTOR_USE_MOCK_DATA") == "false" {
		log.Println("COLLECTOR_USE_MOCK_DATA=false: external provider integration placeholder active, using empty set")
		products = []sourceProduct{}
	}

	rate := 1
	if v := os.Getenv("COLLECTOR_RATE_LIMIT_RPS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			rate = n
		}
	}
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	for _, p := range products {
		<-ticker.C
		pid, err := s.UpsertProduct(ctx, store.Product{
			Name: p.Name, Description: p.Tagline, Category: p.Category,
			Source: "product_hunt", SourceExternalID: p.ExternalID,
			Upvotes: p.Upvotes, CompetitorCount: p.CompetitorCount,
		})
		if err != nil {
			log.Printf("upsert product failed: %v", err)
			continue
		}

		demand := float64(p.Upvotes)*0.6 + float64(p.CompetitorCount)*0.4
		competition := models.CompetitionLevelFromCompetitors(int32(p.CompetitorCount))
		for _, idea := range generateIdeas(p.Name) {
			err := s.UpsertIdea(ctx, store.Idea{IdeaText: idea, SourceProductID: pid, DemandScore: demand, CompetitionLevel: competition})
			if err != nil {
				log.Printf("upsert idea failed: %v", err)
			}
		}
	}

	log.Printf("collector run complete: processed=%d", len(products))
}

func generateIdeas(productName string) []string {
	base := strings.TrimSpace(productName)
	out := make([]string, 0, len(niches))
	for _, n := range niches {
		out = append(out, base+" for "+n)
	}
	return out
}

func mockProducts() []sourceProduct {
	return []sourceProduct{
		{Name: "AI meeting assistant", Tagline: "Summarize meetings", Category: "Productivity", Upvotes: 120, ExternalID: "ph_1", CompetitorCount: 22},
		{Name: "Invoice automation", Tagline: "Auto-payables workflows", Category: "Finance", Upvotes: 78, ExternalID: "ph_2", CompetitorCount: 11},
	}
}
