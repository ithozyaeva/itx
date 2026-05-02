-- Дополнительные флаги пушей для геймификации.
-- Все по умолчанию TRUE, потому что новые пуши осторожные:
-- утренний 1/день, вечерний 1/день только если есть незавершённое,
-- streak — на пересечении порога, raffle-win — только победителю.
ALTER TABLE notification_settings ADD COLUMN IF NOT EXISTS daily_morning  BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE notification_settings ADD COLUMN IF NOT EXISTS daily_evening  BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE notification_settings ADD COLUMN IF NOT EXISTS daily_streak   BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE notification_settings ADD COLUMN IF NOT EXISTS daily_raffle   BOOLEAN NOT NULL DEFAULT TRUE;
