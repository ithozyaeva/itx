CREATE TABLE IF NOT EXISTS "notification_settings" (
  "id" SERIAL PRIMARY KEY,
  "member_id" INTEGER NOT NULL,
  "new_events" BOOLEAN NOT NULL DEFAULT TRUE,
  "remind_week" BOOLEAN NOT NULL DEFAULT TRUE,
  "remind_day" BOOLEAN NOT NULL DEFAULT TRUE,
  "remind_hour" BOOLEAN NOT NULL DEFAULT TRUE,
  "event_start" BOOLEAN NOT NULL DEFAULT TRUE,
  "event_updates" BOOLEAN NOT NULL DEFAULT TRUE,
  "event_cancelled" BOOLEAN NOT NULL DEFAULT TRUE,
  UNIQUE(member_id)
);

CREATE INDEX IF NOT EXISTS "idx_notification_settings_member_id"
  ON "notification_settings"("member_id");

ALTER TABLE "notification_settings"
ADD FOREIGN KEY("member_id") REFERENCES "members"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;
