CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    category TEXT,
    source TEXT NOT NULL,
    source_external_id TEXT,
    upvotes INT NOT NULL DEFAULT 0,
    competitor_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ideas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idea_text TEXT NOT NULL,
    source_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    demand_score NUMERIC(10,2) NOT NULL,
    competition_level TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS products_source_external_id_uniq
    ON products(source, source_external_id);

CREATE UNIQUE INDEX IF NOT EXISTS ideas_source_product_id_idea_text_uniq
    ON ideas(source_product_id, idea_text);

CREATE INDEX IF NOT EXISTS products_category_idx
    ON products(category);

CREATE INDEX IF NOT EXISTS ideas_demand_score_desc_idx
    ON ideas(demand_score DESC);

CREATE INDEX IF NOT EXISTS ideas_competition_level_idx
    ON ideas(competition_level);
