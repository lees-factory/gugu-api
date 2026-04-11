CREATE TABLE IF NOT EXISTS gugu.price_alert (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_user(id),
    sku_id TEXT NOT NULL REFERENCES gugu.sku(id),
    channel TEXT NOT NULL DEFAULT 'EMAIL',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, sku_id)
);

CREATE INDEX IF NOT EXISTS idx_price_alert_user_id ON gugu.price_alert(user_id);
CREATE INDEX IF NOT EXISTS idx_price_alert_sku_id ON gugu.price_alert(sku_id);
