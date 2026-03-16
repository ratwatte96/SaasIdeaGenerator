-- name: UpsertProduct :one
INSERT INTO products (name, description, category, source, source_external_id, upvotes, competitor_count)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (source, source_external_id) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    upvotes = EXCLUDED.upvotes,
    competitor_count = EXCLUDED.competitor_count,
    updated_at = now()
RETURNING *;

-- name: ListProducts :many
SELECT *
FROM products
WHERE ($1::text = '' OR category = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
