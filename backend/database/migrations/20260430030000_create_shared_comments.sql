-- Универсальный модуль комментариев. Сейчас только AI-материалы их имели,
-- расширяем на events; впредь любая сущность подключается через
-- entity_type+entity_id без новой таблицы.
--
-- Миграция атомарно: создаём comments/comment_likes, переносим данные из
-- ai_material_comments/ai_material_comment_likes (FK сохраняем по id),
-- обновляем sequence, дропаем старые таблицы. Старые триггеры на
-- ai_material_comments снимаются автоматически при DROP TABLE CASCADE.

-- 1. Целевые таблицы.

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(32) NOT NULL,
    entity_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    likes_count INTEGER NOT NULL DEFAULT 0,
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_comments_entity ON comments(entity_type, entity_id, created_at);
CREATE INDEX IF NOT EXISTS idx_comments_author ON comments(author_id);

CREATE TABLE IF NOT EXISTS comment_likes (
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (comment_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_comment_likes_member ON comment_likes(member_id);

-- 2. Поле comments_count для events. У AI-материалов уже есть из миграции
--    20260430000000.

ALTER TABLE events
    ADD COLUMN IF NOT EXISTS comments_count INTEGER NOT NULL DEFAULT 0;

-- 3. Перенос существующих комментов и лайков AI-материалов.
--    ai_material_comments.id → comments.id (без сдвига), чтобы внешние ссылки,
--    если они есть в логах/нотификациях/кешах фронта, продолжали резолвиться.

INSERT INTO comments (id, entity_type, entity_id, author_id, body, likes_count, is_hidden, created_at, updated_at)
SELECT id, 'ai_material', material_id, author_id, body, COALESCE(likes_count, 0), is_hidden, created_at, updated_at
FROM ai_material_comments;

INSERT INTO comment_likes (comment_id, member_id, created_at)
SELECT comment_id, member_id, created_at
FROM ai_material_comment_likes;

-- 4. Sequence comments.id_seq нужно подкрутить, иначе следующий BIGSERIAL
--    столкнётся с уже мигрированным id и упадёт по PK.

SELECT setval(
    pg_get_serial_sequence('comments', 'id'),
    COALESCE((SELECT MAX(id) FROM comments), 0) + 1,
    false
);

-- 5. Удаляем старые таблицы вместе с их триггерами/функциями.

DROP TABLE IF EXISTS ai_material_comment_likes;
DROP TABLE IF EXISTS ai_material_comments;

DROP FUNCTION IF EXISTS ai_materials_comments_count_recalc(BIGINT);
DROP FUNCTION IF EXISTS ai_materials_comments_count_after_change();
DROP FUNCTION IF EXISTS ai_material_comments_likes_count_inc();
DROP FUNCTION IF EXISTS ai_material_comments_likes_count_dec();

-- 6. Триггеры новой схемы.

-- 6a. Лайки на комментариях → инкремент/декремент comments.likes_count.

CREATE OR REPLACE FUNCTION comments_likes_count_inc() RETURNS TRIGGER AS $$
BEGIN
    UPDATE comments SET likes_count = likes_count + 1 WHERE id = NEW.comment_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION comments_likes_count_dec() RETURNS TRIGGER AS $$
BEGIN
    UPDATE comments SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = OLD.comment_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_comment_likes_inc ON comment_likes;
CREATE TRIGGER trg_comment_likes_inc
    AFTER INSERT ON comment_likes
    FOR EACH ROW EXECUTE FUNCTION comments_likes_count_inc();

DROP TRIGGER IF EXISTS trg_comment_likes_dec ON comment_likes;
CREATE TRIGGER trg_comment_likes_dec
    AFTER DELETE ON comment_likes
    FOR EACH ROW EXECUTE FUNCTION comments_likes_count_dec();

-- 6b. comments_count на parent — INC/DEC через TG_OP вместо COUNT(*).
--    COUNT(*) на каждое изменение даёт O(N) на популярном треде; INC/DEC
--    держится константой и согласован по стилю с triggers на likes_count.
--    Любой новый entity_type подключается дополнением одной CASE-ветки.

CREATE OR REPLACE FUNCTION comments_count_apply_delta(p_entity_type VARCHAR, p_entity_id BIGINT, p_delta INTEGER) RETURNS VOID AS $$
BEGIN
    IF p_delta = 0 THEN
        RETURN;
    END IF;
    IF p_entity_type = 'ai_material' THEN
        UPDATE ai_materials
            SET comments_count = GREATEST(comments_count + p_delta, 0)
            WHERE id = p_entity_id;
    ELSIF p_entity_type = 'event' THEN
        UPDATE events
            SET comments_count = GREATEST(comments_count + p_delta, 0)
            WHERE id = p_entity_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION comments_count_after_change() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.is_hidden = FALSE THEN
            PERFORM comments_count_apply_delta(NEW.entity_type, NEW.entity_id, 1);
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.is_hidden = FALSE THEN
            PERFORM comments_count_apply_delta(OLD.entity_type, OLD.entity_id, -1);
        END IF;
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Триггер срабатывает только на UPDATE OF is_hidden (см. CREATE TRIGGER ниже),
        -- поэтому достаточно проверить смену значения.
        IF OLD.is_hidden = FALSE AND NEW.is_hidden = TRUE THEN
            PERFORM comments_count_apply_delta(NEW.entity_type, NEW.entity_id, -1);
        ELSIF OLD.is_hidden = TRUE AND NEW.is_hidden = FALSE THEN
            PERFORM comments_count_apply_delta(NEW.entity_type, NEW.entity_id, 1);
        END IF;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_comments_recalc ON comments;
CREATE TRIGGER trg_comments_recalc
    AFTER INSERT OR UPDATE OF is_hidden OR DELETE ON comments
    FOR EACH ROW EXECUTE FUNCTION comments_count_after_change();

-- 7. Восстанавливаем corrected counts на parent после миграции данных
--    (старые ai_materials.comments_count могли разойтись с фактом).
--    Делается ровно один раз здесь, дальше живёт через INC/DEC триггер.

UPDATE ai_materials SET comments_count = (
    SELECT COUNT(*) FROM comments
    WHERE entity_type = 'ai_material' AND entity_id = ai_materials.id AND is_hidden = FALSE
);
