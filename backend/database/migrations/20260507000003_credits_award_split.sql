-- Разделение награды инвайтеру на «первая покупка» и «recurring»:
--   - первая покупка реферала любым способом (Boosty-anchor или credits) →
--     крупная единоразовая выплата (default 50% от рублёвой цены).
--   - каждый месяц активной подписки реферала → recurring-выплата
--     (default 20%), идемпотентность по period_key=YYYY-MM в source_type.
--
-- Reason-константа `referral_purchase` из v1 сменилась на две —
-- `referral_purchase_first` и `referral_purchase_recurring`. Никаких
-- транзакций со старым reason ещё нет (миграция v1 запускается в этом
-- же релизе), поэтому переименовывать существующие данные не требуется.

DROP INDEX IF EXISTS rct_idem_idx;

CREATE UNIQUE INDEX rct_idem_idx
    ON referral_credit_transactions (member_id, reason, source_type, source_id)
    WHERE reason IN (
        'referal_conversion',
        'referral_purchase_first',
        'referral_purchase_recurring'
    );

-- Новые параметры. share старого `referral_purchase_share` подняли до 0.2
-- (recurring-выплата), для первой покупки — отдельный ключ.
UPDATE app_settings SET value = '0.2', updated_at = now()
    WHERE key = 'referral_purchase_share';

INSERT INTO app_settings(key, value) VALUES
    ('referral_first_purchase_share', '0.5')
ON CONFLICT (key) DO NOTHING;
