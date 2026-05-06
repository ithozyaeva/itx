-- price_credits — цена тарифа в реферальных кредитах. NULL = тариф нельзя
-- купить за credits (только Boosty). Стартовое значение: 1 credit ≈ 1 ₽
-- (price_cents/100), админ подкорректирует руками.
ALTER TABLE subscription_tiers
    ADD COLUMN IF NOT EXISTS price_credits INTEGER NULL;

UPDATE subscription_tiers
    SET price_credits = price_cents / 100
    WHERE price_credits IS NULL AND price_cents IS NOT NULL;

-- manual_tier_expires_at — срок действия купленного тарифа. NULL =
-- бессрочно (как у админского override через /suboverride). Истечение
-- проверяется в PeriodicCheck/CheckAndSyncUser перед расчётом
-- EffectiveTierID; просроченный manual сбрасывается, тир падает на
-- resolved (по anchor-чату).
ALTER TABLE subscription_users
    ADD COLUMN IF NOT EXISTS manual_tier_expires_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS su_manual_expires_idx
    ON subscription_users (manual_tier_expires_at)
    WHERE manual_tier_expires_at IS NOT NULL;
