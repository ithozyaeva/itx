-- Лайки на комментарии к AI-материалам. Структура зеркальная к
-- ai_material_likes: денормализованный счётчик в ai_material_comments,
-- триггеры поддерживают синхрон.

ALTER TABLE ai_material_comments
    ADD COLUMN IF NOT EXISTS likes_count INTEGER NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS ai_material_comment_likes (
    comment_id BIGINT NOT NULL REFERENCES ai_material_comments(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (comment_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_ai_material_comment_likes_member
    ON ai_material_comment_likes(member_id);

CREATE OR REPLACE FUNCTION ai_material_comments_likes_count_inc() RETURNS TRIGGER AS $$
BEGIN
    UPDATE ai_material_comments SET likes_count = likes_count + 1 WHERE id = NEW.comment_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ai_material_comments_likes_count_dec() RETURNS TRIGGER AS $$
BEGIN
    UPDATE ai_material_comments SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = OLD.comment_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_ai_material_comment_likes_inc ON ai_material_comment_likes;
CREATE TRIGGER trg_ai_material_comment_likes_inc
    AFTER INSERT ON ai_material_comment_likes
    FOR EACH ROW EXECUTE FUNCTION ai_material_comments_likes_count_inc();

DROP TRIGGER IF EXISTS trg_ai_material_comment_likes_dec ON ai_material_comment_likes;
CREATE TRIGGER trg_ai_material_comment_likes_dec
    AFTER DELETE ON ai_material_comment_likes
    FOR EACH ROW EXECUTE FUNCTION ai_material_comments_likes_count_dec();
