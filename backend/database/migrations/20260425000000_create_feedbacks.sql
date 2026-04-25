CREATE TABLE IF NOT EXISTS feedbacks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES members(id) ON DELETE SET NULL,
    score INT NOT NULL CHECK (score >= 0 AND score <= 10),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_feedbacks_created_at ON feedbacks(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_feedbacks_user_id ON feedbacks(user_id);

INSERT INTO permissions (name)
SELECT name FROM (VALUES ('can_view_admin_feedback')) AS new_perms(name)
WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.name = new_perms.name);

INSERT INTO role_permissions (role, permission_id)
SELECT 'ADMIN', p.id
FROM permissions p
WHERE p.name = 'can_view_admin_feedback'
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp WHERE rp.role = 'ADMIN' AND rp.permission_id = p.id
  );
