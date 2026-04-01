-- name: CreateSKUPriceHistory :exec
INSERT INTO gugu.sku_price_history (
    sku_id,
    price,
    currency,
    recorded_at,
    change_value
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: ListSKUPriceHistoriesBySKUID :many
SELECT
    sku_id,
    price,
    currency,
    recorded_at,
    change_value
FROM gugu.sku_price_history
WHERE sku_id = $1 AND currency = $2
ORDER BY recorded_at DESC;
