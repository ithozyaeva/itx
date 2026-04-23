-- Импортирует в subscription_chats все активные tracked_chats, которых
-- там ещё нет. Эти чаты когда-то попали в tracked (для учёта активности),
-- но не были заведены в систему подписок. Их нельзя «перерегистрировать»
-- через handleMyChatMemberUpdated — Telegram не шлёт my_chat_member при
-- повторном добавлении бота, если он уже в чате.
--
-- После миграции чаты видны в админке «Подписки → Чаты» с фильтром
-- «Без привязки». Тиры/категории/приоритеты админ проставляет руками.
--
-- Фильтр tc.chat_id < -1000000000000 отсекает битые legacy-id basic group
-- (сейчас такой один: -5203153204), которые не являются supergroup и
-- не могут работать с invite-link через Bot API.

INSERT INTO subscription_chats (id, title, chat_type, priority)
SELECT tc.chat_id, tc.title, tc.chat_type, 0
FROM tracked_chats tc
LEFT JOIN subscription_chats sc ON sc.id = tc.chat_id
WHERE tc.is_active
  AND sc.id IS NULL
  AND tc.chat_id < -1000000000000;
