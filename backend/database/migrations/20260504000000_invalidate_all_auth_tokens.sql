-- Массовая инвалидация всех Telegram-токенов в связи с уязвимостью
-- /api/auth/telegram/refresh: ранее эндпоинт принимал base64(telegram_id) и
-- возвращал действующий токен жертвы, что позволяло выписать чужую сессию.
--
-- Здесь не удаляем строки, чтобы при следующем логине через бот
-- CreateOrUpdateToken просто перезаписал token + expired_at без вставки
-- дубликата (telegram_id — UNIQUE). Сдвигаем expired_at в прошлое: middleware
-- вернёт 401, и refresh-хендлер тоже отвергнет такой токен.
UPDATE "auth_tokens"
SET "expired_at" = NOW() - INTERVAL '1 day';
