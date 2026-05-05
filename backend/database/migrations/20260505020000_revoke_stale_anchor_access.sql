-- Чистим стейл записи в subscription_user_chat_access, где chat_id указывает
-- на anchor-чат. Anchor-чаты определяют тир, не являются объектом доступа,
-- и не должны попадать в subscription_user_chat_access. До исправления
-- логики revoke (которое пропускает anchors) periodic-check мог занести их
-- туда; после фикса записи остаются висеть с revoked_at IS NULL и не
-- пересматриваются ни одним loop'ом — гасим явно одним апдейтом.
UPDATE subscription_user_chat_access
SET revoked_at = NOW()
WHERE revoked_at IS NULL
  AND chat_id IN (
    SELECT id FROM subscription_chats WHERE anchor_for_tier_id IS NOT NULL
  );
