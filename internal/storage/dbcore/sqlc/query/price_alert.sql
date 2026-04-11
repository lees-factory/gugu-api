-- name: CreatePriceAlert :exec
INSERT INTO gugu.price_alert (
    id, user_id, sku_id, channel, enabled, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: FindPriceAlertByUserIDAndSKUID :one
SELECT id, user_id, sku_id, channel, enabled, created_at
FROM gugu.price_alert
WHERE user_id = $1 AND sku_id = $2
ORDER BY created_at DESC, id DESC
LIMIT 1;

-- name: ListPriceAlertsBySKUID :many
SELECT id, user_id, sku_id, channel, enabled, created_at
FROM gugu.price_alert
WHERE sku_id = $1 AND enabled = TRUE;

-- name: ListPriceAlertsByProductID :many
SELECT pa.id, pa.user_id, pa.sku_id, pa.channel, pa.enabled, pa.created_at
FROM gugu.price_alert pa
JOIN gugu.sku s ON s.id = pa.sku_id
WHERE s.product_id = $1 AND pa.enabled = TRUE;

-- name: ListPriceAlertsByProductIDs :many
SELECT pa.id, pa.user_id, pa.sku_id, pa.channel, pa.enabled, pa.created_at
FROM gugu.price_alert pa
JOIN gugu.sku s ON s.id = pa.sku_id
WHERE s.product_id = ANY($1::text[]) AND pa.enabled = TRUE;

-- name: ListPriceAlertsByUserID :many
SELECT id, user_id, sku_id, channel, enabled, created_at
FROM gugu.price_alert
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdatePriceAlertEnabled :execrows
UPDATE gugu.price_alert
SET enabled = $2
WHERE id = $1;

-- name: UpdatePriceAlertSettings :execrows
UPDATE gugu.price_alert
SET channel = $2,
    enabled = $3
WHERE id = $1;
