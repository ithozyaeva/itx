-- Снимаем username с дубликатов: оставляем за тем, у кого больший id
-- (последняя по времени запись — её username с большей вероятностью совпадает
-- с актуальным Telegram-username; UPDATE-ы старых строк id не меняют).
UPDATE members SET username = ''
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (PARTITION BY LOWER(username) ORDER BY id DESC) AS rn
        FROM members
        WHERE username <> ''
    ) t
    WHERE rn > 1
);

-- Partial UNIQUE: пустые username не учитываются (у Telegram-аккаунтов без @username
-- поле остаётся пустым, и таких может быть много).
CREATE UNIQUE INDEX IF NOT EXISTS members_username_unique
    ON members (LOWER(username))
    WHERE username <> '';
