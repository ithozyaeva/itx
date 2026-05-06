-- referral_credit_transactions — отдельная валюта «реферальных кредитов»,
-- не смешанная с игровыми point_transactions. Тратится только на покупку
-- подписки (subscription_purchase). Начисляется за реф-конверсию,
-- покупку подписки приглашённым пользователем и вручную админом.

CREATE TABLE referral_credit_transactions (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    reason VARCHAR(50) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    source_id BIGINT NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX rct_member_idx
    ON referral_credit_transactions (member_id, created_at DESC);

-- Идемпотентность для двух reason'ов: повторный referal_conversion для
-- одной (link_id) или referral_purchase для одной (purchase_tx_id) — no-op.
CREATE UNIQUE INDEX rct_idem_idx
    ON referral_credit_transactions (member_id, reason, source_type, source_id)
    WHERE reason IN ('referal_conversion', 'referral_purchase');
