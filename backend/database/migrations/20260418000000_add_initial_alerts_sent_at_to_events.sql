-- Флаг отправки инициализирующих алертов.
-- Нужен, чтобы при APP_MODE=api (бот в отдельном процессе) бот на NL мог
-- подхватить события без начального алерта и доотправить их.
ALTER TABLE events ADD COLUMN IF NOT EXISTS initial_alerts_sent_at TIMESTAMP NULL;

-- Backfill: прошедшие события не должны пытаться отправить алерт.
UPDATE events SET initial_alerts_sent_at = CURRENT_TIMESTAMP
WHERE initial_alerts_sent_at IS NULL AND date < CURRENT_TIMESTAMP;

-- Backfill: будущие события, у которых уже есть подписки — значит, алерт уже уехал раньше.
UPDATE events SET initial_alerts_sent_at = CURRENT_TIMESTAMP
WHERE initial_alerts_sent_at IS NULL
  AND EXISTS (SELECT 1 FROM event_alert_subscriptions WHERE event_id = events.id);
