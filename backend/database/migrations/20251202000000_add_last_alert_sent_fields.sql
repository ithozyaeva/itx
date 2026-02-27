ALTER TABLE "events" 
ADD COLUMN IF NOT EXISTS "last_repeating_alert_sent_at" TIMESTAMP NULL;
