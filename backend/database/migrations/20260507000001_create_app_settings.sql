-- app_settings — конфигурируемые из админки параметры приложения,
-- которые меняются без деплоя. JSONB-значения, чтобы один key мог хранить
-- скаляр, объект или массив без миграций схемы.

CREATE TABLE app_settings (
    key VARCHAR(100) PRIMARY KEY,
    value JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO app_settings(key, value) VALUES
    ('referral_purchase_share', '0.1'),         -- 10% от рублёвой цены идёт инвайтеру при покупке реферала
    ('referral_conversion_credits', '30'),      -- credits за саму конверсию (преемственность с прежними 30 points)
    ('subscription_purchase_days', '30'),       -- срок продления тарифа за credits, в днях
    ('points_to_credits_rate', 'null')          -- курс конвертации points→credits, включится отдельным PR
ON CONFLICT (key) DO NOTHING;
