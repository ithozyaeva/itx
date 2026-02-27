INSERT INTO permissions (name) VALUES ('can_edit_platform_mentor');

INSERT INTO role_permissions (role, permission_id)
SELECT 'MENTOR', id
FROM permissions
WHERE name = 'can_edit_platform_mentor';