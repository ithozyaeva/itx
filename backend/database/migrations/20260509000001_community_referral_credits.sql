-- Расширение идемпотентного индекса на referral_credit_transactions:
-- добавляем reason 'community_referral' (награда за привлечение нового
-- юзера в сообщество через персональный deeplink, отдельно от legacy
-- 'referal_conversion' который начисляется за вакансии).
--
-- Новый source_type='community_referral', source_id=referee_member_id —
-- гарантирует «один award за пару (referrer, referee) на всю жизнь».

DROP INDEX IF EXISTS rct_idem_idx;

CREATE UNIQUE INDEX rct_idem_idx
    ON referral_credit_transactions (member_id, reason, source_type, source_id)
    WHERE reason IN (
        'referal_conversion',
        'community_referral',
        'referral_purchase_first',
        'referral_purchase_recurring'
    );
