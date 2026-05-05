-- Флаг «доступ выдан вручную». Когда админ-человек добавляет юзера в
-- content-чат через нативный Telegram «Add member» (не через invite-link
-- бота), это видно в chat_member событии (update.from != user). Бот
-- проставляет is_manual=true; CheckAndSyncUser и DryRunCheckUser пропускают
-- такие записи в revoke-loop, чтобы периодик не вышибал юзеров, которым
-- админ выдал доступ за заслуги вне их подписочного тира.
--
-- DEFAULT FALSE для обратной совместимости со всеми существующими
-- записями — они продолжают трактоваться как обычный auto-grant.
ALTER TABLE subscription_user_chat_access
ADD COLUMN IF NOT EXISTS is_manual BOOLEAN NOT NULL DEFAULT FALSE;

-- Индекс полу-частичный: revoke-loop фильтрует по revoked_at IS NULL,
-- быстрая отсечка manual=true в этом срезе.
CREATE INDEX IF NOT EXISTS idx_sub_access_manual_active
ON subscription_user_chat_access(user_id, chat_id)
WHERE revoked_at IS NULL AND is_manual = TRUE;
