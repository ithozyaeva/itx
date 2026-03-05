-- Event recordings
ALTER TABLE events ADD COLUMN IF NOT EXISTS recording_url TEXT DEFAULT '';

-- Kudos (wall of thanks)
CREATE TABLE IF NOT EXISTS kudos (
    id BIGSERIAL PRIMARY KEY,
    from_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    to_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT kudos_no_self CHECK (from_id != to_id)
);
CREATE INDEX IF NOT EXISTS idx_kudos_to_id ON kudos(to_id);
CREATE INDEX IF NOT EXISTS idx_kudos_created_at ON kudos(created_at DESC);

-- Seasons
CREATE TABLE IF NOT EXISTS seasons (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    status TEXT DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Raffles
CREATE TABLE IF NOT EXISTS raffles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    prize TEXT NOT NULL,
    ticket_cost INT NOT NULL DEFAULT 10,
    max_tickets INT DEFAULT 0,
    ends_at TIMESTAMPTZ NOT NULL,
    status TEXT DEFAULT 'ACTIVE',
    winner_id BIGINT REFERENCES members(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS raffle_tickets (
    id BIGSERIAL PRIMARY KEY,
    raffle_id BIGINT NOT NULL REFERENCES raffles(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    bought_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_raffle_tickets_raffle ON raffle_tickets(raffle_id);
CREATE INDEX IF NOT EXISTS idx_raffle_tickets_member ON raffle_tickets(member_id);

-- Guilds
CREATE TABLE IF NOT EXISTS guilds (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    icon TEXT DEFAULT 'users',
    color TEXT DEFAULT '#6366f1',
    owner_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS guild_members (
    guild_id BIGINT NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (guild_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_guild_members_member ON guild_members(member_id);

-- Points for kudos and raffle
ALTER TABLE point_transactions DROP CONSTRAINT IF EXISTS point_transactions_member_id_reason_source_type_source_id_key;
