CREATE TABLE IF NOT EXISTS "event_tags" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "event_event_tags" (
  "event_id" INTEGER NOT NULL,
  "event_tag_id" INTEGER NOT NULL,
  PRIMARY KEY (event_id, event_tag_id)
);

ALTER TABLE "event_event_tags"
ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "event_event_tags"
ADD FOREIGN KEY("event_tag_id") REFERENCES "event_tags"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;