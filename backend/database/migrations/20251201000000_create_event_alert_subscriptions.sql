CREATE TABLE IF NOT EXISTS "event_alert_subscriptions" (
  "id" SERIAL PRIMARY KEY,
  "event_id" INTEGER NOT NULL,
  "member_id" INTEGER NOT NULL,
  "status" VARCHAR(50) NOT NULL DEFAULT 'PENDING',
  "reminder_sent_at" TIMESTAMP NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(event_id, member_id)
);

CREATE INDEX IF NOT EXISTS "idx_event_alert_subscriptions_event_id" ON "event_alert_subscriptions"("event_id");
CREATE INDEX IF NOT EXISTS "idx_event_alert_subscriptions_member_id" ON "event_alert_subscriptions"("member_id");
CREATE INDEX IF NOT EXISTS "idx_event_alert_subscriptions_status" ON "event_alert_subscriptions"("status");

ALTER TABLE "event_alert_subscriptions"
ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "event_alert_subscriptions"
ADD FOREIGN KEY("member_id") REFERENCES "members"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;
