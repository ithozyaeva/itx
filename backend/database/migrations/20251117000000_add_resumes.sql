CREATE TABLE IF NOT EXISTS "resumes" (
    "id" SERIAL PRIMARY KEY,
    "tg_id" BIGINT NOT NULL,
    "file_path" VARCHAR(512) NOT NULL,
    "file_name" VARCHAR(255) NOT NULL,
    "work_experience" TEXT,
    "desired_position" VARCHAR(255),
    "work_format" VARCHAR(32),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "members"
    ALTER COLUMN "telegram_id" SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS "members_telegram_id_unique"
    ON "members" ("telegram_id");

ALTER TABLE "resumes"
    ADD CONSTRAINT "resumes_tg_id_fkey"
    FOREIGN KEY ("tg_id")
    REFERENCES "members" ("telegram_id")
    ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS "resumes_tg_id_idx" ON "resumes" ("tg_id");

CREATE INDEX IF NOT EXISTS "resumes_work_format_idx" ON "resumes" ("work_format");

INSERT INTO "permissions" ("name")
SELECT 'can_view_admin_resumes'
WHERE NOT EXISTS (
    SELECT 1 FROM "permissions" WHERE "name" = 'can_view_admin_resumes'
);

INSERT INTO "role_permissions" ("role", "permission_id")
SELECT 'ADMIN', p.id
FROM permissions p
WHERE p.name = 'can_view_admin_resumes'
  AND NOT EXISTS (
        SELECT 1 FROM role_permissions rp
        WHERE rp.role = 'ADMIN' AND rp.permission_id = p.id
    );

