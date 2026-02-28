INSERT INTO permission_models (name) VALUES ('can_view_admin_points'), ('can_edit_admin_points') ON CONFLICT DO NOTHING;
