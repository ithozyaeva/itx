ALTER TABLE "mentorsTags" RENAME TO "mentors_tags";

ALTER TABLE "mentors_tags" RENAME COLUMN "mentorId" TO "mentor_id";
ALTER TABLE "mentors_tags" RENAME COLUMN "tagId" TO "tag_id";

ALTER TABLE "mentors_tags" DROP CONSTRAINT IF EXISTS mentors_tags_mentorid_fkey;
ALTER TABLE "mentors_tags" DROP CONSTRAINT IF EXISTS mentors_tags_tagid_fkey;

ALTER TABLE "mentors_tags"
ADD CONSTRAINT fk_mentor FOREIGN KEY ("mentor_id") REFERENCES "mentors"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "mentors_tags"
ADD CONSTRAINT fk_tag FOREIGN KEY ("tag_id") REFERENCES "profTags"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;