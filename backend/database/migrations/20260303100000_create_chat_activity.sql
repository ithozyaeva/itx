-- Таблица отслеживаемых чатов
CREATE TABLE IF NOT EXISTS "tracked_chats" (
  "id" SERIAL PRIMARY KEY,
  "chat_id" BIGINT NOT NULL UNIQUE,
  "title" VARCHAR(255) NOT NULL,
  "chat_type" VARCHAR(50) NOT NULL DEFAULT 'supergroup',
  "is_active" BOOLEAN NOT NULL DEFAULT true,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица сообщений (одна запись = одно сообщение)
CREATE TABLE IF NOT EXISTS "chat_messages" (
  "id" SERIAL PRIMARY KEY,
  "chat_id" BIGINT NOT NULL,
  "telegram_user_id" BIGINT NOT NULL,
  "telegram_username" VARCHAR(255),
  "telegram_first_name" VARCHAR(255),
  "sent_at" TIMESTAMP NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрых агрегаций
CREATE INDEX IF NOT EXISTS idx_chat_messages_chat_id_sent_at ON chat_messages(chat_id, sent_at);
CREATE INDEX IF NOT EXISTS idx_chat_messages_telegram_user_id ON chat_messages(telegram_user_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_sent_at ON chat_messages(sent_at);

-- Заполняем известные чаты
INSERT INTO tracked_chats (chat_id, title) VALUES
  (-1001847344728, 'IT-ХОЗЯЕВА: Основная беседа'),
  (-1002459179718, 'IT-X: Игры'),
  (-1002252250331, 'IT-X: Новички в IT'),
  (-1002403724947, 'IT-X: Инвестиции, темки, деньги, крипта'),
  (-1002350487206, 'IT-X: Frontend'),
  (-1002270225191, 'IT-X: Backend'),
  (-1002623226260, 'IT-X: Железо и техника'),
  (-1002436169789, 'IT-X: English club'),
  (-1002322657217, 'IT-X: Еда, кулинария, рецепты'),
  (-1002394057716, 'IT-X: DevOps, Infrastructure'),
  (-1002492438891, 'IT-X: Рефералки, поиск работы (оффтоп)'),
  (-1002889736366, 'IT-X: Хобби и творчество'),
  (-1002863279888, 'IT-X: Психология, менталка'),
  (-1002598592247, 'IT-X: ИИ, Нейросети, AI'),
  (-1002544203048, 'IT-X: Разбор резюме'),
  (-1002615578774, 'IT-X: Архитектурный'),
  (-1003822754192, 'IT-X: Вакансии, Работа'),
  (-1003340480893, 'IT-X: Рефералки, поиск работы'),
  (-1002936591048, 'IT-X: Музыка')
ON CONFLICT (chat_id) DO NOTHING;
