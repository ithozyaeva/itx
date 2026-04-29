-- Расширение subscription_tiers: цены, ссылки на оплату, публичное описание
ALTER TABLE subscription_tiers
    ADD COLUMN IF NOT EXISTS price_cents INTEGER,
    ADD COLUMN IF NOT EXISTS boosty_url VARCHAR(255),
    ADD COLUMN IF NOT EXISTS public_description TEXT,
    ADD COLUMN IF NOT EXISTS features JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS is_public BOOLEAN NOT NULL DEFAULT FALSE;

-- Новичок (259 ₽) — boosty URL пока неизвестен, скрыт из публичного списка
UPDATE subscription_tiers SET
    name = 'Новичок',
    price_cents = 25900,
    public_description = 'Вход в ИТ-сообщество',
    features = '["Тематические ИТ-чаты","Спортивная беседа","Книжный клуб","Встречи по английскому"]'::jsonb,
    is_public = FALSE
WHERE slug = 'beginner';

-- Бригадир (520 ₽)
UPDATE subscription_tiers SET
    name = 'Бригадир',
    price_cents = 52000,
    boosty_url = 'https://boosty.to/jointime/purchase/3150816',
    public_description = 'Старт в ИТ-сообществе',
    features = '["Все возможности «Новичка»","Доступ к рефералкам","Чат с инвестициями"]'::jsonb,
    is_public = TRUE
WHERE slug = 'foreman';

-- Хозяин (2000 ₽, на Boosty первый месяц со скидкой 50%)
UPDATE subscription_tiers SET
    name = 'Хозяин',
    price_cents = 200000,
    boosty_url = 'https://boosty.to/jointime/purchase/3150814',
    public_description = 'Закрытый клуб и живые беседы',
    features = '["Все возможности «Бригадира»","Закрытый клуб","Еженедельные встречи в формате живой беседы","Личная благодарность"]'::jsonb,
    is_public = TRUE
WHERE slug = 'master';

-- Мастер (5200 ₽) — был KING на лендинге
UPDATE subscription_tiers SET
    name = 'Мастер',
    price_cents = 520000,
    boosty_url = 'https://boosty.to/jointime/purchase/967625',
    public_description = 'Персональное менторство и продвижение',
    features = '["Все возможности «Хозяина»","Реклама ресурсов","Секция с менторингом на сайте","Верхняя позиция в таблице менторов","Личная консультация по карьере или технологиям"]'::jsonb,
    is_public = TRUE
WHERE slug = 'king';
