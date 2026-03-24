-- Ретроактивно начисляет ачивку "Заполнить профиль полностью" (profile_complete, 20 баллов)
-- тем участникам, у кого заполнены first_name, last_name, bio и birthday,
-- но транзакция ещё не существует (ON CONFLICT DO NOTHING — идемпотентно).
INSERT INTO point_transactions (member_id, amount, reason, source_type, source_id, description, created_at)
SELECT
    m.id,
    20,
    'profile_complete',
    'profile',
    m.id,
    'Полностью заполненный профиль',
    NOW()
FROM members m
WHERE m.first_name  != ''
  AND m.last_name   != ''
  AND m.bio         != ''
  AND m.birthday    IS NOT NULL
ON CONFLICT DO NOTHING;
