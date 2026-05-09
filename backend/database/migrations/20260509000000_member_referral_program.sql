-- Реферальная программа на подписку IT-X (отдельно от referal_links для вакансий).
--
-- referral_code — персональный код юзера для приглашения в сообщество.
--   8 символов из алфавита Crockford-base32 (без 0/O/1/I/L/U → 32 символа,
--   ~10^12 комбинаций, не угадать перебором). NULL до backfill, после —
--   обязательно через application logic. Уникальный индекс — гарантия attribution.
--
-- referred_by_member_id — кто пригласил (members.id). FK с ON DELETE SET NULL —
--   если инвайтера удалили, атрибуция теряется, но Боб остаётся валидным юзером.
--
-- referral_welcome_seen_at — когда Боб закрыл welcome-баннер. NULL = ещё не видел.

ALTER TABLE members
    ADD COLUMN IF NOT EXISTS referral_code VARCHAR(16) NULL,
    ADD COLUMN IF NOT EXISTS referred_by_member_id BIGINT NULL REFERENCES members(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS referral_welcome_seen_at TIMESTAMPTZ NULL;

-- Уникальность кода. Не enforced до backfill (разрешаем NULL в существующих
-- строках до того как Go-стартап сгенерирует код всем). После того как все
-- строки заполнены, можно сделать NOT NULL отдельной миграцией.
CREATE UNIQUE INDEX IF NOT EXISTS members_referral_code_idx
    ON members (referral_code)
    WHERE referral_code IS NOT NULL;

CREATE INDEX IF NOT EXISTS members_referred_by_idx
    ON members (referred_by_member_id)
    WHERE referred_by_member_id IS NOT NULL;
