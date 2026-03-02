-- Expire all auth tokens created before today (2026-03-02).
-- Tokens are valid for 1 month, so tokens created today have expired_at >= 2026-04-01.
-- Setting expired_at to NOW() forces re-login on next request.
UPDATE auth_tokens
SET expired_at = NOW()
WHERE expired_at < NOW() + INTERVAL '29 days';
