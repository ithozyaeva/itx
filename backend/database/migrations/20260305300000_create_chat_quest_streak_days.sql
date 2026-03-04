CREATE TABLE IF NOT EXISTS chat_quest_streak_days (
    id BIGSERIAL PRIMARY KEY,
    quest_id BIGINT NOT NULL REFERENCES chat_quests(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL,
    day DATE NOT NULL,
    UNIQUE(quest_id, member_id, day)
);

CREATE INDEX idx_chat_quest_streak_days_member ON chat_quest_streak_days(quest_id, member_id, day);
