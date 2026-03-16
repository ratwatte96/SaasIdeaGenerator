package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Store struct{ db *sql.DB }

func New(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) ListIdeas(ctx context.Context, category, competition string, minDemand float64, limit, offset int) ([]Idea, error) {
	query := `
SELECT i.id, i.idea_text, i.source_product_id, i.demand_score, i.competition_level, i.created_at,
       p.name, p.category
FROM ideas i
JOIN products p ON p.id = i.source_product_id
WHERE ($1 = '' OR p.category = $1)
  AND ($2 = '' OR i.competition_level = $2)
  AND i.demand_score >= $3
ORDER BY i.demand_score DESC, i.created_at DESC
LIMIT $4 OFFSET $5`
	rows, err := s.db.QueryContext(ctx, query, category, competition, minDemand, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ideas := []Idea{}
	for rows.Next() {
		var i Idea
		if err := rows.Scan(&i.ID, &i.IdeaText, &i.SourceProductID, &i.DemandScore, &i.CompetitionLevel, &i.CreatedAt, &i.ProductName, &i.ProductCategory); err != nil {
			return nil, err
		}
		ideas = append(ideas, i)
	}
	return ideas, rows.Err()
}

func (s *Store) GetIdea(ctx context.Context, id string) (Idea, error) {
	query := `SELECT i.id, i.idea_text, i.source_product_id, i.demand_score, i.competition_level, i.created_at, p.name, p.category
FROM ideas i JOIN products p ON p.id=i.source_product_id WHERE i.id=$1`
	var i Idea
	err := s.db.QueryRowContext(ctx, query, id).Scan(&i.ID, &i.IdeaText, &i.SourceProductID, &i.DemandScore, &i.CompetitionLevel, &i.CreatedAt, &i.ProductName, &i.ProductCategory)
	return i, err
}

func (s *Store) RelatedIdeas(ctx context.Context, sourceProductID, excludeID string, limit int) ([]Idea, error) {
	query := `SELECT i.id, i.idea_text, i.source_product_id, i.demand_score, i.competition_level, i.created_at, p.name, p.category
FROM ideas i JOIN products p ON p.id=i.source_product_id
WHERE i.source_product_id=$1 AND i.id<>$2
ORDER BY i.demand_score DESC LIMIT $3`
	rows, err := s.db.QueryContext(ctx, query, sourceProductID, excludeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Idea{}
	for rows.Next() {
		var i Idea
		if err := rows.Scan(&i.ID, &i.IdeaText, &i.SourceProductID, &i.DemandScore, &i.CompetitionLevel, &i.CreatedAt, &i.ProductName, &i.ProductCategory); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

func (s *Store) ListProducts(ctx context.Context, category string, limit, offset int) ([]Product, error) {
	query := `SELECT id,name,description,category,source,COALESCE(source_external_id,''),upvotes,competitor_count,created_at
FROM products WHERE ($1='' OR category=$1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := s.db.QueryContext(ctx, query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Category, &p.Source, &p.SourceExternalID, &p.Upvotes, &p.CompetitorCount, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) UpsertProduct(ctx context.Context, p Product) (string, error) {
	q := `INSERT INTO products (name,description,category,source,source_external_id,upvotes,competitor_count)
VALUES ($1,$2,$3,$4,$5,$6,$7)
ON CONFLICT (source,source_external_id) DO UPDATE
SET name=EXCLUDED.name, description=EXCLUDED.description, category=EXCLUDED.category,
    upvotes=EXCLUDED.upvotes, competitor_count=EXCLUDED.competitor_count, updated_at=now()
RETURNING id`
	var id string
	err := s.db.QueryRowContext(ctx, q, p.Name, p.Description, p.Category, p.Source, p.SourceExternalID, p.Upvotes, p.CompetitorCount).Scan(&id)
	return id, err
}

func (s *Store) UpsertIdea(ctx context.Context, i Idea) error {
	q := `INSERT INTO ideas (idea_text,source_product_id,demand_score,competition_level)
VALUES ($1,$2,$3,$4)
ON CONFLICT (source_product_id, idea_text) DO UPDATE
SET demand_score=EXCLUDED.demand_score, competition_level=EXCLUDED.competition_level`
	_, err := s.db.ExecContext(ctx, q, strings.TrimSpace(i.IdeaText), i.SourceProductID, i.DemandScore, i.CompetitionLevel)
	return err
}

func ParseLimitOffset(limitS, offsetS string) (int, int) {
	limit, offset := 50, 0
	if _, err := fmt.Sscan(limitS, &limit); err != nil || limit <= 0 || limit > 200 {
		limit = 50
	}
	if _, err := fmt.Sscan(offsetS, &offset); err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset
}
