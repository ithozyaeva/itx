-- Seed-пул дейликов. Состав согласован с пользователем: 17 задач,
-- из них 5 случайных в день. Можно править через админку (PR следующих
-- этапов), но стартовая версия задаётся миграцией, чтобы не плодить
-- ручную работу при первичной настройке.
--
-- ON CONFLICT (code) — миграция идемпотентна: повторный прогон не плодит
-- дубликаты и не затирает изменения админа в admin-only полях, кроме
-- дефолтов из этого seed-а.

INSERT INTO daily_tasks (code, title, description, icon, tier, points, target, trigger_key, active) VALUES
-- Engagement (10 баллов)
('view_ai_material', 'Открой свежий AI-материал', 'Загляни в раздел AI-материалов', 'sparkles', 'engagement', 10, 1, 'view_ai_material', TRUE),
('view_events',      'Заглянуть в афишу событий', 'Посмотри, что нас ждёт в ближайшее время', 'calendar', 'engagement', 10, 1, 'view_events', TRUE),
('view_leaderboard', 'Посмотреть рейтинг',         'Узнай, кто на верхушке этой недели', 'trophy', 'engagement', 10, 1, 'view_leaderboard', TRUE),
('view_mentor',      'Открой профиль ментора',     'Полистай карточку любого ментора', 'graduation-cap', 'engagement', 10, 1, 'view_mentor', TRUE),
('view_kudos',       'Загляни на стену благодарностей', 'Посмотри, кого сообщество благодарит сегодня', 'heart', 'engagement', 10, 1, 'view_kudos', TRUE),

-- Light Action (15 баллов)
('like_comment',     'Лайкни чужой комментарий', 'Поддержи кого-нибудь лайком', 'thumbs-up', 'light', 15, 1, 'like_comment', TRUE),
('send_kudos',       'Отправь благодарность',    'Скажи спасибо кому-то из сообщества', 'hand-heart', 'light', 15, 1, 'send_kudos', TRUE),
('register_event',   'Зарегистрируйся на событие', 'Запишись на любое предстоящее мероприятие', 'calendar-plus', 'light', 15, 1, 'register_event', TRUE),
('view_ai_2',        'Прочти 2 AI-материала',    'Двойная доза знаний', 'book-open', 'light', 15, 2, 'view_ai_material', TRUE),
('bookmark_ai',      'Сохрани материал в закладки', 'Сохрани полезный материал, чтобы не потерять', 'bookmark', 'light', 15, 1, 'bookmark_ai', TRUE),

-- Meaningful Action (20 баллов)
('post_comment',     'Напиши комментарий',       'Поделись мыслью под любым материалом', 'message-square', 'meaningful', 20, 1, 'post_comment', TRUE),
('chat_3_msgs',      'Напиши 3 сообщения в чате', 'Поучаствуй в обсуждениях сообщества', 'messages-square', 'meaningful', 20, 3, 'chat_message', TRUE),
('update_profile',   'Обнови поле в профиле',    'Допиши био, обнови контакты или теги', 'user-cog', 'meaningful', 20, 1, 'update_profile', TRUE),
('create_referal',   'Создай реф-ссылку',        'Поделись с миром свежей вакансией', 'link-2', 'meaningful', 20, 1, 'create_referal', TRUE),
('buy_raffle_ticket','Купи билет в розыгрыш',    'Поучаствуй в обычном розыгрыше', 'ticket', 'meaningful', 20, 1, 'buy_raffle_ticket', TRUE),
('leave_review',     'Оставь отзыв',             'Поделись впечатлениями о сообществе или услуге', 'star', 'meaningful', 20, 1, 'leave_review', TRUE),

-- Big Action (30 баллов)
('attend_event',     'Отметься на прошедшем событии', 'Подтверди участие в недавнем эфире', 'check-circle-2', 'big', 30, 1, 'attend_event', TRUE),
('update_resume',    'Обнови резюме',            'Освежи свою карточку специалиста', 'file-text', 'big', 30, 1, 'update_resume', TRUE)
ON CONFLICT (code) DO NOTHING;
