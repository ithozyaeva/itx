ALTER TABLE "members" 
    ADD COLUMN "telegram_id" BIGINT, 
    ADD COLUMN "first_name" VARCHAR(255), 
    ADD COLUMN "last_name" VARCHAR(255), 
    ADD COLUMN "username" VARCHAR(255), 
    ADD COLUMN "role" VARCHAR(255) DEFAULT 'UNSUBSCRIBER';

ALTER TABLE "members" 
    DROP COLUMN "tg", 
    DROP COLUMN "name";

CREATE TABLE IF NOT EXISTS "auth_tokens" (
    "id" SERIAL PRIMARY KEY,
    "telegram_id" BIGINT UNIQUE NOT NULL,
    "expired_at" TIMESTAMP,
    "token" VARCHAR(255)
); 
