CREATE TABLE IF NOT EXISTS "chat_quests" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "quest_type" VARCHAR(50) NOT NULL DEFAULT 'message_count',
  "chat_id" BIGINT,
  "target_count" INT NOT NULL,
  "points_reward" INT NOT NULL DEFAULT 10,
  "starts_at" TIMESTAMP NOT NULL,
  "ends_at" TIMESTAMP NOT NULL,
  "is_active" BOOLEAN NOT NULL DEFAULT true,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "chat_quest_progress" (
  "id" SERIAL PRIMARY KEY,
  "quest_id" INT NOT NULL REFERENCES chat_quests(id) ON DELETE CASCADE,
  "member_id" BIGINT NOT NULL,
  "current_count" INT NOT NULL DEFAULT 0,
  "completed" BOOLEAN NOT NULL DEFAULT false,
  "completed_at" TIMESTAMP,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(quest_id, member_id)
);

CREATE INDEX idx_chat_quest_progress_member ON chat_quest_progress(member_id);
CREATE INDEX idx_chat_quest_progress_quest ON chat_quest_progress(quest_id, completed);
CREATE INDEX idx_chat_quests_active ON chat_quests(is_active, starts_at, ends_at);
