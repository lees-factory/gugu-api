-- price_alert를 사용자+SKU 단일 알림 모델로 정규화한다.
-- 1) 기존 (user_id, sku_id, channel) 중복 데이터를 (user_id, sku_id) 1건으로 정리
-- 2) unique key를 (user_id, sku_id)로 전환

WITH ranked AS (
    SELECT
        id,
        user_id,
        sku_id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, sku_id
            ORDER BY enabled DESC, created_at DESC, id DESC
        ) AS rn
    FROM gugu.price_alert
),
to_delete AS (
    SELECT id
    FROM ranked
    WHERE rn > 1
)
DELETE FROM gugu.price_alert
WHERE id IN (SELECT id FROM to_delete);

ALTER TABLE gugu.price_alert
    DROP CONSTRAINT IF EXISTS price_alert_user_id_product_id_channel_key;

ALTER TABLE gugu.price_alert
    DROP CONSTRAINT IF EXISTS price_alert_user_id_sku_id_channel_key;

ALTER TABLE gugu.price_alert
    DROP CONSTRAINT IF EXISTS price_alert_user_id_sku_id_key;

ALTER TABLE gugu.price_alert
    ADD CONSTRAINT price_alert_user_id_sku_id_key
        UNIQUE (user_id, sku_id);
