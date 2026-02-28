CREATE TABLE IF NOT EXISTS referral_conversions (
    id SERIAL PRIMARY KEY,
    referral_link_id INTEGER NOT NULL REFERENCES referal_links(id) ON DELETE CASCADE,
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    converted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT referral_conversions_unique UNIQUE(referral_link_id, member_id)
);

CREATE INDEX IF NOT EXISTS idx_referral_conversions_link_id ON referral_conversions(referral_link_id);
