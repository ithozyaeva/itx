-- referee attribution: при /start ref_<id> в боте Боб сохраняется в
-- Redis-pending. Когда auth-handler создаёт ему members-запись, link-id
-- переезжает сюда и фиксируется навсегда. Без FK к referal_links —
-- ссылка может быть удалена/заморожена, но факт «пришёл от X» остаётся.
ALTER TABLE members
    ADD COLUMN IF NOT EXISTS referred_by_link_id BIGINT NULL,
    ADD COLUMN IF NOT EXISTS referral_welcome_seen_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS members_referred_by_link_idx
    ON members (referred_by_link_id)
    WHERE referred_by_link_id IS NOT NULL;
