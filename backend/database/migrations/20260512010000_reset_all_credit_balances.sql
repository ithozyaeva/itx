-- Security reset: обнуляем ВСЕ положительные балансы referral_credit_transactions.
--
-- Причина: до этой миграции findReferrerForMember в subscription.go имел
-- fallback на legacy referral_conversions → referal_links атрибуцию. В
-- результате покупки подписки пользователями, прошедшими когда-то по
-- чужим job-referral-ссылкам, давали авторам тех ссылок community-кредиты
-- (referral_purchase_first / referral_purchase_recurring), хотя
-- референты НЕ были атрибутированы через личный deeplink (referred_by_member_id).
-- Это смешивало две разные продуктовые сущности и приводило к «бесплатным»
-- кредитам без явного community-приглашения.
--
-- Кодовый фикс (subscription.go: удалён fallback) предотвращает новые
-- неправильные начисления. Эта миграция приводит баланс всех юзеров к
-- нулю одним admin_manual-сбросом, чтобы перейти на корректную атрибуцию
-- с чистого листа. Дальше начисления идут только через явный
-- referred_by_member_id (community-deeplink ref_<code> в боте).
--
-- Идемпотентность: source_type='security_reset' — уникальный маркер.
-- NOT EXISTS guard защищает от повторного запуска. Дополнительная защита
-- поверх дедупа миграций по filename (database/db.go).
--
-- Безопасность: операция в транзакции (db.go обёртывает каждую миграцию
-- в BEGIN/COMMIT). Миграции применяются на старте бэка ДО app.Listen() —
-- конкурентные Spend (FOR UPDATE) невозможны в момент применения.

WITH balances AS (
    SELECT member_id, SUM(amount) AS balance
    FROM referral_credit_transactions
    GROUP BY member_id
)
INSERT INTO referral_credit_transactions (
    member_id, amount, reason, source_type, source_id, description
)
SELECT
    member_id,
    -balance,
    'admin_manual',
    'security_reset',
    0,
    'Сброс баланса: переход на строгую community-only атрибуцию рефералов'
FROM balances
WHERE balance > 0
  AND NOT EXISTS (
      SELECT 1 FROM referral_credit_transactions wo
      WHERE wo.member_id = balances.member_id
        AND wo.source_type = 'security_reset'
  );
