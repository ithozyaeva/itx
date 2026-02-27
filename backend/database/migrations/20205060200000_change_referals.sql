DROP TABLE "referal_links_tags";


CREATE TABLE IF NOT EXISTS "referal_links_tags" (
  "referal_link_id" INTEGER NOT NULL,
  "prof_tag_id" INTEGER NOT NULL,
  PRIMARY KEY ("referal_link_id", "prof_tag_id")
);

ALTER TABLE "referal_links_tags"
ADD FOREIGN KEY("referal_link_id") REFERENCES "referal_links"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "referal_links_tags"
ADD FOREIGN KEY("prof_tag_id") REFERENCES "profTags"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;