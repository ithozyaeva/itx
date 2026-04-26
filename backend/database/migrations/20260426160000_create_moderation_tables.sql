-- Сохраняем telegram message_id у трекаемых сообщений, чтобы /cleanup мог
-- их удалять через Telegram API. Старые записи останутся с NULL — их удалить
-- невозможно (Telegram API требует message_id).
ALTER TABLE chat_messages
    ADD COLUMN IF NOT EXISTS telegram_message_id INTEGER;

CREATE INDEX IF NOT EXISTS idx_chat_messages_chat_user_sent
    ON chat_messages (chat_id, telegram_user_id, sent_at);

-- Аудит модерационных действий: ban, unban, mute, unmute, cleanup, voteban_mute.
CREATE TABLE IF NOT EXISTS bot_moderation_actions (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    target_user_id BIGINT NOT NULL,
    actor_user_id BIGINT NOT NULL,
    action VARCHAR(32) NOT NULL,
    reason TEXT,
    duration_seconds INTEGER,
    expires_at TIMESTAMPTZ,
    meta JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bot_moderation_actions_chat_user
    ON bot_moderation_actions (chat_id, target_user_id);
CREATE INDEX IF NOT EXISTS idx_bot_moderation_actions_created_at
    ON bot_moderation_actions (created_at);

-- Голосования за бан (voteban) — по одному открытому на (chat_id, target).
CREATE TABLE IF NOT EXISTS bot_votebans (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    chat_title VARCHAR(255) NOT NULL DEFAULT '',
    target_user_id BIGINT NOT NULL,
    target_username VARCHAR(255) NOT NULL DEFAULT '',
    target_first_name VARCHAR(255) NOT NULL DEFAULT '',
    initiator_user_id BIGINT NOT NULL,
    trigger_message_id INTEGER,
    poll_message_id INTEGER NOT NULL,
    required_votes INTEGER NOT NULL DEFAULT 5,
    mute_seconds INTEGER NOT NULL DEFAULT 3600,
    expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    finalized_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_bot_votebans_open
    ON bot_votebans (chat_id, target_user_id)
    WHERE status = 'open';
CREATE INDEX IF NOT EXISTS idx_bot_votebans_open_expires
    ON bot_votebans (expires_at)
    WHERE status = 'open';

-- Голоса по voteban: уникальные на (voteban_id, voter).
CREATE TABLE IF NOT EXISTS bot_voteban_votes (
    voteban_id BIGINT NOT NULL REFERENCES bot_votebans(id) ON DELETE CASCADE,
    voter_user_id BIGINT NOT NULL,
    vote SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (voteban_id, voter_user_id)
);
