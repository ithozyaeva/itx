CREATE TABLE IF NOT EXISTS "chat_highlights" (
  "id" SERIAL PRIMARY KEY,
  "chat_id" BIGINT NOT NULL,
  "message_id" INTEGER NOT NULL,
  "author_telegram_id" BIGINT NOT NULL,
  "author_username" VARCHAR(255) DEFAULT '',
  "author_first_name" VARCHAR(255) DEFAULT '',
  "message_text" TEXT NOT NULL,
  "highlighted_by" BIGINT NOT NULL,
  "member_id" INTEGER,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_chat_highlights_created_at"
  ON "chat_highlights"("created_at" DESC);

ALTER TABLE "chat_highlights"
ADD FOREIGN KEY("member_id") REFERENCES "members"("id")
ON UPDATE NO ACTION ON DELETE SET NULL;
