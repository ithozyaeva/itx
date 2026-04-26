-- Глобальные баны: блокировка пользователя во всех чатах сразу + отказ
-- в выдаче доступа по подписке. Снимается через /globalunban.
CREATE TABLE IF NOT EXISTS bot_global_bans (
    user_id BIGINT PRIMARY KEY,
    banned_by BIGINT NOT NULL,
    reason TEXT,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bot_global_bans_expires
    ON bot_global_bans (expires_at);
