-- name: ListPriceHistoriesByProductID :many
SELECT
    product_id,
    price,
    currency,
    recorded_at,
    change_value
FROM gugu.product_price_histories
WHERE product_id = $1
ORDER BY recorded_at DESC;
