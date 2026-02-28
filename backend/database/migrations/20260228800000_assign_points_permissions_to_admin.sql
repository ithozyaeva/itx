-- Ensure points permissions exist
INSERT INTO permissions (name)
SELECT name FROM (VALUES ('can_view_admin_points'), ('can_edit_admin_points')) AS new_perms(name)
WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.name = new_perms.name);

-- Assign to ADMIN role
INSERT INTO role_permissions (role, permission_id)
SELECT 'ADMIN', p.id
FROM permissions p
WHERE p.name IN ('can_view_admin_points', 'can_edit_admin_points')
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp WHERE rp.role = 'ADMIN' AND rp.permission_id = p.id
  );
