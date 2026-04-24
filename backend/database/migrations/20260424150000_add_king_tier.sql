INSERT INTO subscription_tiers (slug, name, level) VALUES
    ('king', 'King', 4)
ON CONFLICT (slug) DO NOTHING;
