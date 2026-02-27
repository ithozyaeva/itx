CREATE TABLE IF NOT EXISTS "referal_links" (
  "id" SERIAL PRIMARY KEY,
  "author_id" INTEGER NOT NULL,
  "company" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "grade" VARCHAR NOT NULL,
  "status" VARCHAR NOT NULL,
  "vacations_count" INT DEFAULT 1 NULL
);

CREATE TABLE IF NOT EXISTS "referal_links_tags" (
  "id" SERIAL PRIMARY KEY,
  "referal_link_id" INTEGER NOT NULL,
  "prof_tag_id" INTEGER NOT NULL
);

ALTER TABLE "referal_links"
ADD FOREIGN KEY("author_id") REFERENCES "members"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "referal_links_tags"
ADD FOREIGN KEY("referal_link_id") REFERENCES "referal_links"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "referal_links_tags"
ADD FOREIGN KEY("prof_tag_id") REFERENCES "profTags"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;