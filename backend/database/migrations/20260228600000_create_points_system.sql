CREATE TABLE IF NOT EXISTS point_transactions (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    amount INT NOT NULL,
    reason VARCHAR(50) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    source_id BIGINT NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_point_transactions_unique
    ON point_transactions (member_id, reason, source_type, source_id);

CREATE INDEX IF NOT EXISTS idx_point_transactions_member_id
    ON point_transactions (member_id);

CREATE INDEX IF NOT EXISTS idx_point_transactions_created_at
    ON point_transactions (created_at);
