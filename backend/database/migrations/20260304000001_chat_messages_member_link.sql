ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS member_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_chat_messages_member_id ON chat_messages(member_id);
