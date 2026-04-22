-- Backfill tracked_chats from subscription_chats.
--
-- Авто-трекинг (PR #260) добавляет чат в tracked_chats, когда бота добавляют
-- в группу. Чаты, которые уже были в subscription_chats до этого PR, в
-- tracked_chats не попали — и для них не работает автоочистка service-
-- сообщений (она фильтрует по IsTrackedChat). Одноразово заливаем все
-- недостающие записи; у новых чатов всё продолжит работать через
-- ChatActivityService.AddTrackedChat.

INSERT INTO tracked_chats (chat_id, title, chat_type, is_active)
SELECT sc.id, sc.title, sc.chat_type, TRUE
FROM subscription_chats sc
WHERE NOT EXISTS (
    SELECT 1 FROM tracked_chats tc WHERE tc.chat_id = sc.id
);
