CREATE TABLE IF NOT EXISTS "events" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR NOT NULL,
  "description" VARCHAR NULL,
  "date" TIMESTAMP NOT NULL,
  "place_type" VARCHAR(255) NOT NULL,
  "custom_place_type" VARCHAR NULL,
  "event_type" VARCHAR(255) NOT NULL,
  "place" VARCHAR NOT NULL,
  "video_link" VARCHAR NULL,
  "open" BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS "event_hosts" (
  "event_id" INTEGER NOT NULL,
  "member_id" INTEGER NOT NULL,
  PRIMARY KEY (event_id, member_id)
);

ALTER TABLE "event_hosts"
ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "event_members" (
  "event_id" INTEGER NOT NULL,
  "member_id" INTEGER NOT NULL,
  PRIMARY KEY (event_id, member_id)
);

ALTER TABLE "event_members"
ADD FOREIGN KEY("event_id") REFERENCES "events"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;
