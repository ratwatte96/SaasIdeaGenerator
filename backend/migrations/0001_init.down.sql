DROP INDEX IF EXISTS ideas_competition_level_idx;
DROP INDEX IF EXISTS ideas_demand_score_desc_idx;
DROP INDEX IF EXISTS products_category_idx;
DROP INDEX IF EXISTS ideas_source_product_id_idea_text_uniq;
DROP INDEX IF EXISTS products_source_external_id_uniq;

DROP TABLE IF EXISTS ideas;
DROP TABLE IF EXISTS products;
