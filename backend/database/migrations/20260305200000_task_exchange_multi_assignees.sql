-- Create assignees junction table
CREATE TABLE IF NOT EXISTS task_exchange_assignees (
    task_id BIGINT NOT NULL REFERENCES task_exchanges(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (task_id, member_id)
);

CREATE INDEX IF NOT EXISTS idx_task_exchange_assignees_member ON task_exchange_assignees(member_id);

-- Add max_assignees column
ALTER TABLE task_exchanges ADD COLUMN IF NOT EXISTS max_assignees INT NOT NULL DEFAULT 1;

-- Migrate existing assignee data
INSERT INTO task_exchange_assignees (task_id, member_id, created_at)
SELECT id, assignee_id, NOW()
FROM task_exchanges
WHERE assignee_id IS NOT NULL
ON CONFLICT DO NOTHING;

-- Drop old column and index
DROP INDEX IF EXISTS idx_task_exchange_assignee;
ALTER TABLE task_exchanges DROP COLUMN IF EXISTS assignee_id;
