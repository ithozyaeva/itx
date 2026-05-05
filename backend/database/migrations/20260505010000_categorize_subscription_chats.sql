-- Раскладывает 37 subscription_chats по логическим категориям с emoji и priority,
-- чтобы /substatus и /mygroups в боте показывали группированный список вместо
-- одной кучи в «Прочее». Категории сортируются по MAX(priority) DESC
-- (см. formatChatsGrouped в backend/internal/bot/subscription.go).
--
-- Идемпотентно: точечный UPDATE по title. Если чат был переименован в TG,
-- его придётся проставить вручную через админку.

-- VIP (priority 100)
UPDATE subscription_chats SET category = 'VIP', emoji = '💎', priority = 100
    WHERE title IN ('AI-X', 'IT-X: База Стародубцева');

-- Главное (priority 95)
UPDATE subscription_chats SET category = 'Главное', emoji = '🏠', priority = 95
    WHERE title = 'IT-ХОЗЯЕВА: Основная беседа';

-- ИИ (priority 85) — выше Разработки
UPDATE subscription_chats SET category = 'ИИ', emoji = '🤖', priority = 85
    WHERE title = 'IT-X: ИИ, Нейросети, AI';

-- Разработка (priority 80)
UPDATE subscription_chats SET category = 'Разработка', emoji = '💻', priority = 80
    WHERE title IN (
        'IT-X: Backend',
        'IT-X: Frontend',
        'IT-X: Golang банда',
        'IT-X: Python',
        'IT-X: DevOps, Infrastructure',
        'IT-X: Архитектурный',
        'IT-X: Алгосы',
        'IT-X: ИБ, кибербеза'
    );

-- Карьера (priority 70)
UPDATE subscription_chats SET category = 'Карьера', emoji = '💼', priority = 70
    WHERE title IN (
        'IT-X: Вакансии, Работа',
        'IT-X: Корпорации. Карьера',
        'IT-X: Новички в IT',
        'IT-X: Разбор резюме',
        'IT-X: Собесы',
        'IT-X: Рефералки, поиск работы',
        'IT-X: Рефералки, поиск работы (оффтоп)'
    );

-- Обучение (priority 60)
UPDATE subscription_chats SET category = 'Обучение', emoji = '📚', priority = 60
    WHERE title IN (
        'IT-X: English club. Английский',
        'IT-X: Книги, литература, чтение',
        'IT-X: Технический книжный клуб'
    );

-- Финансы (priority 55)
UPDATE subscription_chats SET category = 'Финансы', emoji = '💰', priority = 55
    WHERE title = 'IT-X: Инвестиции, темки, деньги, крипта';

-- Города (priority 50)
UPDATE subscription_chats SET category = 'Города', emoji = '🌆', priority = 50
    WHERE title IN ('IT-X: Москва', 'IT-X: Питер', 'IT-X: Воронеж');

-- Хобби и досуг (priority 40)
UPDATE subscription_chats SET category = 'Хобби и досуг', emoji = '🎮', priority = 40
    WHERE title IN (
        'IT-X: Аниме',
        'IT-X: Игры',
        'IT-X: Музыка',
        'IT-X: Хобби и творчество',
        'IT-X: Еда, кулинария, рецепты',
        'IT-X: Путешествия',
        'IT-X: Спорт, здоровье',
        'IT-X: Железо и техника'
    );

-- Менталка (priority 30)
UPDATE subscription_chats SET category = 'Менталка', emoji = '🧠', priority = 30
    WHERE title = 'IT-X: Психология, менталка';

-- Болталка (priority 10)
UPDATE subscription_chats SET category = 'Болталка', emoji = '💬', priority = 10
    WHERE title IN ('IT-X: /b/ - Бред, политика, срачи', 'IT-X: VPN');
