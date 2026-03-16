-- name: UpsertIdea :exec
INSERT INTO ideas (idea_text, source_product_id, demand_score, competition_level)
VALUES ($1, $2, $3, $4)
ON CONFLICT (source_product_id, idea_text) DO UPDATE
SET demand_score = EXCLUDED.demand_score,
    competition_level = EXCLUDED.competition_level;

-- name: ListIdeas :many
SELECT i.*, p.name AS product_name, p.category AS product_category
FROM ideas i
JOIN products p ON p.id = i.source_product_id
WHERE ($1::text = '' OR p.category = $1)
  AND ($2::text = '' OR i.competition_level = $2)
  AND i.demand_score >= $3
ORDER BY i.demand_score DESC
LIMIT $4 OFFSET $5;

-- name: GetIdeaByID :one
SELECT i.*, p.name AS product_name, p.category AS product_category
FROM ideas i
JOIN products p ON p.id = i.source_product_id
WHERE i.id = $1;

-- name: RelatedIdeas :many
SELECT i.*, p.name AS product_name, p.category AS product_category
FROM ideas i
JOIN products p ON p.id = i.source_product_id
WHERE i.source_product_id = $1
  AND i.id <> $2
ORDER BY i.demand_score DESC
LIMIT $3;
