-- expired_notified_at — отметка, что бот уже отправил алерт «срок истёк» по
-- этому модерационному действию. Watcher проверяет expires_at <= NOW() и
-- expired_notified_at IS NULL, шлёт сообщение в чат(ы) и проставляет отметку.
ALTER TABLE bot_moderation_actions
    ADD COLUMN IF NOT EXISTS expired_notified_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_bot_moderation_actions_expiry_pending
    ON bot_moderation_actions (expires_at)
    WHERE expires_at IS NOT NULL AND expired_notified_at IS NULL;
