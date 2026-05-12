-- Списываем исторические referal_conversion-начисления.
--
-- Бизнес-решение: реферальные ссылки на вакансии (legacy) и community-рефералы
-- через личный deeplink — две разные системы. Кредитная валюта остаётся только
-- за community (см. удалённый AwardForConversion). Для консистентности убираем
-- из баланса всё, что было начислено через 'referal_conversion'.
--
-- Стратегия: для каждого юзера, у которого был хоть один положительный
-- referal_conversion, создаём ОДНУ компенсирующую admin_manual-запись. Сумма —
-- LEAST(заработано_по_referal_conversion, текущий_баланс): не уводим юзера в
-- минус, если он уже потратил часть кредитов на подписку (Spend держит
-- инвариант balance >= 0, отрицательный баланс ломает следующее списание).
--
-- Идемпотентность: source_type='referal_conversion_writeoff' — уникальный
-- маркер. Повторный прогон через NOT EXISTS — no-op. Дополнительная защита
-- поверх дедупа миграций по filename (database/db.go).

WITH user_stats AS (
    SELECT
        member_id,
        SUM(CASE WHEN reason = 'referal_conversion' AND amount > 0 THEN amount ELSE 0 END) AS conversion_earned,
        SUM(amount) AS current_balance
    FROM referral_credit_transactions
    GROUP BY member_id
)
INSERT INTO referral_credit_transactions (
    member_id, amount, reason, source_type, source_id, description
)
SELECT
    member_id,
    -LEAST(conversion_earned, current_balance),
    'admin_manual',
    'referal_conversion_writeoff',
    0,
    'Корректировка: реферальные ссылки на вакансии больше не дают кредитов'
FROM user_stats
WHERE conversion_earned > 0
  AND current_balance > 0
  AND NOT EXISTS (
      SELECT 1 FROM referral_credit_transactions wo
      WHERE wo.member_id = user_stats.member_id
        AND wo.source_type = 'referal_conversion_writeoff'
  );
