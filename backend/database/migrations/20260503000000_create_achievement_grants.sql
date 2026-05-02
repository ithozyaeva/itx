-- Явные ачивки, выдаваемые не по count-of-reason (как в AllAchievements),
-- а по конкретному событию — например, при completion челленджа с
-- achievement_code в шаблоне.
CREATE TABLE IF NOT EXISTS achievement_grants (
    member_id  BIGINT      NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    code       VARCHAR(64) NOT NULL,
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (member_id, code)
);
CREATE INDEX IF NOT EXISTS idx_achievement_grants_member ON achievement_grants(member_id);
