-- name: CreatePriceHistory :exec
INSERT INTO gugu.product_price_history (
    product_id,
    price,
    currency,
    recorded_at,
    change_value
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: ListPriceHistoriesByProductID :many
SELECT
    product_id,
    price,
    currency,
    recorded_at,
    change_value
FROM gugu.product_price_history
WHERE product_id = $1
ORDER BY recorded_at DESC;
