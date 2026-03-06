-- Добавляем новые чаты в аналитику
INSERT INTO tracked_chats (chat_id, title) VALUES
  (-1002552216248, 'IT-X: База Стародубцева'),
  (-1003229483670, 'IT-X: /b/ - Бред, политика, срачи'),
  (-1002264228346, 'IT-X: VPN')
ON CONFLICT (chat_id) DO NOTHING;
