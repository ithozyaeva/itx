-- Gamification foundation: dailies, check-ins, streaks, challenges + raffle kind extension.
-- Behavior is implemented in subsequent PRs; this migration creates the schema only.

-- Pool of daily task templates (admin-managed)
CREATE TABLE IF NOT EXISTS daily_tasks (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    icon VARCHAR(64) NOT NULL DEFAULT 'circle',
    tier VARCHAR(16) NOT NULL,
    points INT NOT NULL,
    target INT NOT NULL DEFAULT 1,
    trigger_key VARCHAR(64) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_daily_tasks_trigger_key ON daily_tasks(trigger_key) WHERE active;
CREATE INDEX IF NOT EXISTS idx_daily_tasks_tier ON daily_tasks(tier) WHERE active;

-- One row per MSK-day with the chosen task ids (shared by all members)
CREATE TABLE IF NOT EXISTS daily_task_sets (
    day DATE PRIMARY KEY,
    task_ids BIGINT[] NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Per-member progress on a given day's task
CREATE TABLE IF NOT EXISTS daily_task_progress (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    day DATE NOT NULL,
    task_id BIGINT NOT NULL REFERENCES daily_tasks(id) ON DELETE CASCADE,
    progress INT NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ,
    awarded BOOLEAN NOT NULL DEFAULT FALSE,
    bonus_awarded BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (member_id, day, task_id)
);
CREATE INDEX IF NOT EXISTS idx_dtp_member_day ON daily_task_progress(member_id, day);

-- Daily check-ins (1 per member per MSK-day)
CREATE TABLE IF NOT EXISTS daily_check_ins (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    day DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (member_id, day)
);
CREATE INDEX IF NOT EXISTS idx_check_ins_day ON daily_check_ins(day);

-- Per-member streak summary
CREATE TABLE IF NOT EXISTS member_streaks (
    member_id BIGINT PRIMARY KEY REFERENCES members(id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    last_check_in_date DATE,
    freezes_available INT NOT NULL DEFAULT 1,
    freeze_week_year INT,
    freeze_week_num INT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Challenge templates (admin-managed)
CREATE TABLE IF NOT EXISTS challenge_templates (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    icon VARCHAR(64) NOT NULL DEFAULT 'trophy',
    kind VARCHAR(16) NOT NULL,
    metric_key VARCHAR(64) NOT NULL,
    target INT NOT NULL,
    reward_points INT NOT NULL,
    achievement_code VARCHAR(64),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_challenge_templates_kind ON challenge_templates(kind) WHERE active;

-- Spawned instances of a template for a specific period (week / month)
CREATE TABLE IF NOT EXISTS challenge_instances (
    id BIGSERIAL PRIMARY KEY,
    template_id BIGINT NOT NULL REFERENCES challenge_templates(id) ON DELETE CASCADE,
    kind VARCHAR(16) NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    period_key VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (template_id, period_key)
);
CREATE INDEX IF NOT EXISTS idx_ci_active ON challenge_instances(starts_at, ends_at);
CREATE INDEX IF NOT EXISTS idx_ci_period ON challenge_instances(period_key);

-- Per-member progress on a specific challenge instance
CREATE TABLE IF NOT EXISTS challenge_progress (
    id BIGSERIAL PRIMARY KEY,
    instance_id BIGINT NOT NULL REFERENCES challenge_instances(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    progress INT NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ,
    awarded BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (instance_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_cp_member ON challenge_progress(member_id);

-- Raffle extension: differentiate auto-generated daily raffles from manual ones
ALTER TABLE raffles ADD COLUMN IF NOT EXISTS kind VARCHAR(16) NOT NULL DEFAULT 'manual';
ALTER TABLE raffles ADD COLUMN IF NOT EXISTS entry_rule VARCHAR(24) NOT NULL DEFAULT 'purchase';
ALTER TABLE raffles ADD COLUMN IF NOT EXISTS day_key DATE;
CREATE UNIQUE INDEX IF NOT EXISTS uniq_raffle_day_key ON raffles(day_key) WHERE kind = 'daily';
