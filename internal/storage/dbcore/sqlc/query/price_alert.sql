-- name: CreatePriceAlert :exec
INSERT INTO gugu.price_alert (
    id, user_id, product_id, channel, enabled, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: FindPriceAlertByUserIDAndProductID :one
SELECT id, user_id, product_id, channel, enabled, created_at
FROM gugu.price_alert
WHERE user_id = $1 AND product_id = $2;

-- name: ListPriceAlertsByProductID :many
SELECT id, user_id, product_id, channel, enabled, created_at
FROM gugu.price_alert
WHERE product_id = $1 AND enabled = TRUE;

-- name: ListPriceAlertsByProductIDs :many
SELECT id, user_id, product_id, channel, enabled, created_at
FROM gugu.price_alert
WHERE product_id = ANY($1::text[]) AND enabled = TRUE;

-- name: ListPriceAlertsByUserID :many
SELECT id, user_id, product_id, channel, enabled, created_at
FROM gugu.price_alert
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdatePriceAlertEnabled :execrows
UPDATE gugu.price_alert
SET enabled = $2
WHERE id = $1;

-- name: DeletePriceAlertByUserIDAndProductID :execrows
DELETE FROM gugu.price_alert
WHERE user_id = $1 AND product_id = $2;
