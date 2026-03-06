-- Добавляем поле exclusive_chat_id для эксклюзивных событий
ALTER TABLE events ADD COLUMN IF NOT EXISTS exclusive_chat_id BIGINT DEFAULT NULL;
ALTER TABLE events ADD COLUMN IF NOT EXISTS exclusive_chat_title VARCHAR(255) DEFAULT '';
