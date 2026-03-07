CREATE TABLE IF NOT EXISTS casino_bets (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id),
    game VARCHAR(20) NOT NULL,
    bet_amount INT NOT NULL,
    multiplier DOUBLE PRECISION NOT NULL DEFAULT 0,
    payout INT NOT NULL DEFAULT 0,
    result JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_casino_bets_member_id ON casino_bets(member_id);
CREATE INDEX idx_casino_bets_created_at ON casino_bets(created_at);
