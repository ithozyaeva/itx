CREATE TABLE IF NOT EXISTS "member_roles" (
    "member_id" INTEGER NOT NULL,
    "role" VARCHAR(255) NOT NULL,
    CONSTRAINT member_roles_unique_member_role UNIQUE(member_id, role),
    CONSTRAINT member_roles_member_fk FOREIGN KEY(member_id) REFERENCES "members"(id) ON DELETE CASCADE
);

INSERT INTO member_roles (member_id, role)
SELECT id, role
FROM members
WHERE role IS NOT NULL;

INSERT INTO member_roles (member_id, role)
SELECT id, 'SUBSCRIBER'
FROM members
WHERE role = 'MENTOR';

ALTER TABLE "members"
DROP COLUMN IF EXISTS "role";

CREATE TABLE IF NOT EXISTS "permissions" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS "role_permissions" (
    role VARCHAR(255) NOT NULL,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role, permission_id)
);

INSERT INTO permissions (name) VALUES
('can_view_admin_panel'),
('can_view_admin_members'),
('can_view_admin_mentors'),
('can_view_admin_events'),
('can_edit_admin_members'),
('can_edit_admin_mentors'),
('can_edit_admin_events'),
('can_view_admin_reviews'),
('can_edit_admin_reviews'),
('can_approved_admin_reviews'),
('can_view_admin_mentors_review'),
('can_edit_admin_mentors_review'),
('can_approve_admin_mentors_review');

INSERT INTO role_permissions (role, permission_id)
SELECT 'ADMIN', id
FROM permissions
WHERE name IN (
    'can_view_admin_panel',
    'can_view_admin_members',
    'can_view_admin_mentors',
    'can_view_admin_events',
    'can_edit_admin_members',
    'can_edit_admin_mentors',
    'can_edit_admin_events',
    'can_view_admin_reviews',
    'can_edit_admin_reviews',
    'can_approved_admin_reviews',
    'can_view_admin_mentors_review',
    'can_edit_admin_mentors_review',
    'can_approve_admin_mentors_review'
);

INSERT INTO role_permissions (role, permission_id)
SELECT 'EVENT_MAKER', id
FROM permissions
WHERE name IN (
    'can_view_admin_panel',
    'can_view_admin_members',
    'can_view_admin_events',
    'can_edit_admin_events'
);
