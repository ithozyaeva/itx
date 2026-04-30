-- Перф-индексы для AI-материалов:
-- 1. pg_trgm GIN на title/summary — ILIKE '%q%' начинает использовать индекс
--    вместо seq-scan, что важно при росте каталога.
-- 2. B-tree на likes_count DESC — ускоряет sort=popular в листинге.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_ai_materials_title_trgm
    ON ai_materials USING GIN (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_ai_materials_summary_trgm
    ON ai_materials USING GIN (summary gin_trgm_ops);

-- Sort=popular: ORDER BY likes_count DESC, created_at DESC.
-- Composite index с второй колонкой как tiebreaker.
CREATE INDEX IF NOT EXISTS idx_ai_materials_popular
    ON ai_materials (likes_count DESC, created_at DESC)
    WHERE is_hidden = FALSE;
