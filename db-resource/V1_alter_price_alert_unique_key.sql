-- price_alert: SKU 기반 unique key 추가
ALTER TABLE gugu.price_alert
    DROP CONSTRAINT IF EXISTS price_alert_user_id_product_id_channel_key;

ALTER TABLE gugu.price_alert
    DROP CONSTRAINT IF EXISTS price_alert_user_id_sku_id_channel_key;

ALTER TABLE gugu.price_alert
    DROP CONSTRAINT IF EXISTS price_alert_user_id_sku_id_key;

ALTER TABLE gugu.price_alert
    ADD CONSTRAINT price_alert_user_id_sku_id_key
        UNIQUE (user_id, sku_id);

CREATE INDEX IF NOT EXISTS idx_price_alert_sku_id ON gugu.price_alert(sku_id);
