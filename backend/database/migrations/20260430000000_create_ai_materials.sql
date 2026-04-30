-- База AI-материалов. Подписчики master+ кладут промты, ссылки на внешние
-- ресурсы и конфиги готовых AI-агентов; листинг + соцфункции (лайки, закладки,
-- комментарии). Гейтинг доступа реализован в middleware на уровне routes,
-- здесь только структура данных.

CREATE TABLE IF NOT EXISTS ai_materials (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    title VARCHAR(120) NOT NULL,
    summary TEXT NOT NULL,
    -- prompt | link | agent — определяет, какой из контентных полей содержит данные
    content_type VARCHAR(16) NOT NULL CHECK (content_type IN ('prompt', 'link', 'agent')),
    -- prompt | skill | library | tutorial | agent — верхнеуровневая категория для фильтрации
    material_kind VARCHAR(16) NOT NULL CHECK (material_kind IN ('prompt', 'skill', 'library', 'tutorial', 'agent')),
    prompt_body TEXT,
    external_url VARCHAR(2048),
    agent_config TEXT,
    likes_count INTEGER NOT NULL DEFAULT 0,
    bookmarks_count INTEGER NOT NULL DEFAULT 0,
    comments_count INTEGER NOT NULL DEFAULT 0,
    -- Мягкое удаление админом (автор удаляет жёстко). Скрытые невидимы в листинге,
    -- но не удаляются — чтобы можно было восстановить или разобрать кейс.
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_materials_author ON ai_materials(author_id);
CREATE INDEX IF NOT EXISTS idx_ai_materials_kind ON ai_materials(material_kind);
CREATE INDEX IF NOT EXISTS idx_ai_materials_created ON ai_materials(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ai_materials_visible ON ai_materials(is_hidden, created_at DESC);

-- Свободные теги; нормализованы до lowercase в сервисе. Без отдельного
-- словаря — для autocomplete агрегируем по join-таблице.
CREATE TABLE IF NOT EXISTS ai_material_tags (
    material_id BIGINT NOT NULL REFERENCES ai_materials(id) ON DELETE CASCADE,
    tag VARCHAR(40) NOT NULL,
    PRIMARY KEY (material_id, tag)
);
CREATE INDEX IF NOT EXISTS idx_ai_material_tags_tag ON ai_material_tags(tag);

CREATE TABLE IF NOT EXISTS ai_material_likes (
    material_id BIGINT NOT NULL REFERENCES ai_materials(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (material_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_ai_material_likes_member ON ai_material_likes(member_id);

CREATE TABLE IF NOT EXISTS ai_material_bookmarks (
    material_id BIGINT NOT NULL REFERENCES ai_materials(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (material_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_ai_material_bookmarks_member ON ai_material_bookmarks(member_id, created_at DESC);

CREATE TABLE IF NOT EXISTS ai_material_comments (
    id BIGSERIAL PRIMARY KEY,
    material_id BIGINT NOT NULL REFERENCES ai_materials(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ai_material_comments_material ON ai_material_comments(material_id, created_at);

-- Триггеры на счётчики. Альтернатива — пересчитывать в Go, но рассинхрон
-- между транзакциями неизбежен; триггер в той же транзакции что вставка.

CREATE OR REPLACE FUNCTION ai_materials_likes_count_inc() RETURNS TRIGGER AS $$
BEGIN
    UPDATE ai_materials SET likes_count = likes_count + 1 WHERE id = NEW.material_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ai_materials_likes_count_dec() RETURNS TRIGGER AS $$
BEGIN
    UPDATE ai_materials SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = OLD.material_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_ai_material_likes_inc ON ai_material_likes;
CREATE TRIGGER trg_ai_material_likes_inc
    AFTER INSERT ON ai_material_likes
    FOR EACH ROW EXECUTE FUNCTION ai_materials_likes_count_inc();

DROP TRIGGER IF EXISTS trg_ai_material_likes_dec ON ai_material_likes;
CREATE TRIGGER trg_ai_material_likes_dec
    AFTER DELETE ON ai_material_likes
    FOR EACH ROW EXECUTE FUNCTION ai_materials_likes_count_dec();

CREATE OR REPLACE FUNCTION ai_materials_bookmarks_count_inc() RETURNS TRIGGER AS $$
BEGIN
    UPDATE ai_materials SET bookmarks_count = bookmarks_count + 1 WHERE id = NEW.material_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ai_materials_bookmarks_count_dec() RETURNS TRIGGER AS $$
BEGIN
    UPDATE ai_materials SET bookmarks_count = GREATEST(bookmarks_count - 1, 0) WHERE id = OLD.material_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_ai_material_bookmarks_inc ON ai_material_bookmarks;
CREATE TRIGGER trg_ai_material_bookmarks_inc
    AFTER INSERT ON ai_material_bookmarks
    FOR EACH ROW EXECUTE FUNCTION ai_materials_bookmarks_count_inc();

DROP TRIGGER IF EXISTS trg_ai_material_bookmarks_dec ON ai_material_bookmarks;
CREATE TRIGGER trg_ai_material_bookmarks_dec
    AFTER DELETE ON ai_material_bookmarks
    FOR EACH ROW EXECUTE FUNCTION ai_materials_bookmarks_count_dec();

-- Комменты: счётчик растёт только на видимые. При hide/unhide пересчитывается.
CREATE OR REPLACE FUNCTION ai_materials_comments_count_recalc(p_material_id BIGINT) RETURNS VOID AS $$
BEGIN
    UPDATE ai_materials
    SET comments_count = (
        SELECT COUNT(*) FROM ai_material_comments
        WHERE material_id = p_material_id AND is_hidden = FALSE
    )
    WHERE id = p_material_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ai_materials_comments_count_after_change() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM ai_materials_comments_count_recalc(OLD.material_id);
        RETURN OLD;
    ELSE
        PERFORM ai_materials_comments_count_recalc(NEW.material_id);
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_ai_material_comments_recalc ON ai_material_comments;
CREATE TRIGGER trg_ai_material_comments_recalc
    AFTER INSERT OR UPDATE OF is_hidden OR DELETE ON ai_material_comments
    FOR EACH ROW EXECUTE FUNCTION ai_materials_comments_count_after_change();
