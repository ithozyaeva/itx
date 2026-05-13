-- Авто-генерация чат-квестов для оживления тихих чатов.
-- auto_generated отделяет машинные квесты от ручных (админка).
-- notification_posted_at — координация API↔бот: API создаёт запись с NULL,
-- бот (APP_MODE=bot, другой сервер) поллит и постит в чат, выставляя время.
ALTER TABLE chat_quests
  ADD COLUMN IF NOT EXISTS auto_generated BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE chat_quests
  ADD COLUMN IF NOT EXISTS notification_posted_at TIMESTAMP NULL;

-- Идемпотентность: один авто-квест в день на чат. Если cron запустится
-- повторно (рестарт API), INSERT … ON CONFLICT DO NOTHING не задвоит.
CREATE UNIQUE INDEX IF NOT EXISTS uniq_chat_quests_auto_chat_day
  ON chat_quests (chat_id, (DATE(starts_at)))
  WHERE auto_generated = TRUE AND chat_id IS NOT NULL;

-- Ускоряет поллинг бота: SELECT … WHERE auto_generated AND notification_posted_at IS NULL.
CREATE INDEX IF NOT EXISTS idx_chat_quests_auto_unposted
  ON chat_quests (notification_posted_at)
  WHERE auto_generated = TRUE AND notification_posted_at IS NULL;
