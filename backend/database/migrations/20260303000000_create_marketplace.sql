CREATE TABLE IF NOT EXISTS marketplace_items (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    price VARCHAR(100) DEFAULT '',
    city VARCHAR(255) DEFAULT '',
    can_ship BOOLEAN DEFAULT FALSE,
    condition VARCHAR(20) DEFAULT 'USED',
    defects TEXT DEFAULT '',
    package_contents TEXT DEFAULT '',
    contact_telegram VARCHAR(255) DEFAULT '',
    contact_email VARCHAR(255) DEFAULT '',
    contact_phone VARCHAR(255) DEFAULT '',
    image_path TEXT DEFAULT '',
    seller_id BIGINT NOT NULL REFERENCES members(id),
    buyer_id BIGINT REFERENCES members(id),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_items_seller_id ON marketplace_items(seller_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_items_status ON marketplace_items(status);
CREATE INDEX IF NOT EXISTS idx_marketplace_items_created_at ON marketplace_items(created_at);
