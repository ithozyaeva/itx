-- Связываем 3 ключевых ежемесячных шаблона с achievement_code.
-- Если seed уже применён — UPDATE без condition на DO NOTHING достаточно
-- идемпотентен (повторный прогон выставит то же значение).
UPDATE challenge_templates SET achievement_code = 'owner_of_month' WHERE code = 'm_owner';
UPDATE challenge_templates SET achievement_code = 'perfect_month'  WHERE code = 'm_check_in_30';
UPDATE challenge_templates SET achievement_code = 'kudos_legend'   WHERE code = 'm_kudos_legend';
