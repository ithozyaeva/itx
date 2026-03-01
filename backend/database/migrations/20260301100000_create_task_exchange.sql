CREATE TABLE IF NOT EXISTS task_exchanges (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    creator_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    assignee_id BIGINT REFERENCES members(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_task_exchange_status ON task_exchanges(status);
CREATE INDEX IF NOT EXISTS idx_task_exchange_creator ON task_exchanges(creator_id);
CREATE INDEX IF NOT EXISTS idx_task_exchange_assignee ON task_exchanges(assignee_id);

-- Permission for approving/rejecting tasks
INSERT INTO permissions (name)
SELECT name FROM (VALUES ('can_approve_platform_tasks')) AS new_perms(name)
WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.name = new_perms.name);

-- Assign permission to ADMIN role
INSERT INTO role_permissions (role, permission_id)
SELECT 'ADMIN', p.id
FROM permissions p
WHERE p.name IN ('can_approve_platform_tasks')
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp WHERE rp.role = 'ADMIN' AND rp.permission_id = p.id
  );
