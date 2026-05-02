-- Seed-пул шаблонов челленджей. Еженедельные стартуют каждый понедельник
-- (3 случайных из активных), ежемесячный — 1 числа (1 случайный).
-- Идемпотентно: повторный прогон не плодит дубликаты благодаря UNIQUE
-- на code.

INSERT INTO challenge_templates
    (code, title, description, icon, kind, metric_key, target, reward_points, achievement_code, active)
VALUES
    -- Еженедельные
    ('w_event_guest',     'Гость недели',         'Посети 2 события за неделю',                'calendar-check', 'weekly', 'events_attended',          2,  200, NULL, TRUE),
    ('w_kudos_social',    'Социалка',             'Получи 30 благодарностей за неделю',         'hand-heart',     'weekly', 'kudos_received',           30, 250, NULL, TRUE),
    ('w_daily_marathon',  'Дейли-марафонец',      'Выполни все 5 дейликов 5 дней за неделю',    'flame',          'weekly', 'all_dailies_days',         5,  300, NULL, TRUE),
    ('w_streaker',        'Огонь не гаснет',      'Сделай 7 check-in’ов за неделю',             'fire',           'weekly', 'check_ins',                7,  200, NULL, TRUE),
    ('w_commenter',       'Голос комьюнити',      'Напиши 10 комментариев за неделю',           'message-square', 'weekly', 'comments_posted',          10, 200, NULL, TRUE),
    ('w_mentor_seeker',   'Охотник за знаниями',  'Открой профили 5 разных менторов',           'graduation-cap', 'weekly', 'mentor_profiles_viewed',   5,  150, NULL, TRUE),
    ('w_ai_reader',       'AI-читатель',          'Прочти 7 AI-материалов за неделю',           'sparkles',       'weekly', 'ai_materials_read',        7,  200, NULL, TRUE),
    ('w_chatter',         'Болтун недели',        'Напиши 50 сообщений в чатах сообщества',     'messages-square','weekly', 'chat_messages',            50, 250, NULL, TRUE),
    -- Ежемесячные
    ('m_owner',           'Хозяин месяца',        'Заработай 1500 баллов за месяц',             'crown',          'monthly','points_earned',            1500, 1000, NULL, TRUE),
    ('m_check_in_30',     'Идеальный месяц',      'Сделай 28+ check-in’ов за месяц',            'calendar',       'monthly','check_ins',                28, 1500, NULL, TRUE),
    ('m_event_pro',       'Завсегдатай',          'Посети 6 событий за месяц',                  'medal',          'monthly','events_attended',          6,  800, NULL, TRUE),
    ('m_kudos_legend',    'Легенда комьюнити',    'Получи 100 благодарностей за месяц',         'heart',          'monthly','kudos_received',           100,1000, NULL, TRUE),
    ('m_dailies_25',      'Дейли-фанат',          'Выполни все 5 дейликов 25 дней за месяц',    'target',         'monthly','all_dailies_days',         25, 1200, NULL, TRUE),
    ('m_referrer',        'Рекрутер',             'Приведи 3 участников по реф-ссылке',         'user-plus',      'monthly','referal_conversions',      3,  1500, NULL, TRUE)
ON CONFLICT (code) DO NOTHING;
