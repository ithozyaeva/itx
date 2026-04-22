-- priority управляет порядком категорий в списках чатов, которые бот шлёт
-- пользователям (/sub, /substatus, /mygroups, рассылки). При группировке по
-- category группы сортируются по MAX(priority) DESC — так «элитные» чаты
-- (База Стародубцева, AI-X) можно вывести наверх, просто проставив им
-- высокий priority; дефолт 0 сохраняет текущее поведение.

ALTER TABLE subscription_chats
    ADD COLUMN IF NOT EXISTS priority INTEGER NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_subscription_chats_priority
    ON subscription_chats (priority DESC);
