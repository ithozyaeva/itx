-- Subscription tiers
CREATE TABLE IF NOT EXISTS subscription_tiers (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    level INT NOT NULL UNIQUE
);

-- Seed default tiers
INSERT INTO subscription_tiers (slug, name, level) VALUES
    ('beginner', 'Beginner', 1),
    ('foreman', 'Foreman', 2),
    ('master', 'Master', 3)
ON CONFLICT (slug) DO NOTHING;

-- Subscription chats (anchor + content)
CREATE TABLE IF NOT EXISTS subscription_chats (
    id BIGINT PRIMARY KEY,
    title VARCHAR(255) NOT NULL DEFAULT '',
    chat_type VARCHAR(50) NOT NULL DEFAULT 'supergroup',
    anchor_for_tier_id INT REFERENCES subscription_tiers(id) ON DELETE SET NULL
);

-- Many-to-many: tier <-> chat (content chats)
CREATE TABLE IF NOT EXISTS subscription_tier_chats (
    tier_id INT NOT NULL REFERENCES subscription_tiers(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL REFERENCES subscription_chats(id) ON DELETE CASCADE,
    PRIMARY KEY (tier_id, chat_id)
);

-- Subscription users (Telegram users tracked for access)
CREATE TABLE IF NOT EXISTS subscription_users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(255),
    full_name VARCHAR(255) NOT NULL DEFAULT '',
    resolved_tier_id INT REFERENCES subscription_tiers(id) ON DELETE SET NULL,
    manual_tier_id INT REFERENCES subscription_tiers(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_check_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User chat access tracking
CREATE TABLE IF NOT EXISTS subscription_user_chat_access (
    user_id BIGINT NOT NULL REFERENCES subscription_users(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL REFERENCES subscription_chats(id) ON DELETE CASCADE,
    granted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    revoked_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, chat_id)
);

-- Subscription audit log
CREATE TABLE IF NOT EXISTS subscription_audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_sub_audit_user_id ON subscription_audit_logs(user_id);
CREATE INDEX idx_sub_audit_created_at ON subscription_audit_logs(created_at);
CREATE INDEX idx_sub_user_chat_access_revoked ON subscription_user_chat_access(revoked_at);
