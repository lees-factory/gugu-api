-- name: UpsertProductPriceSnapshot :exec
INSERT INTO gugu.product_price_snapshot (
    product_id,
    snapshot_date,
    price,
    currency
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (product_id, currency, snapshot_date) DO UPDATE SET
    price = EXCLUDED.price;

-- name: ListProductPriceSnapshotsByDateRange :many
SELECT
    product_id,
    snapshot_date,
    price,
    currency
FROM gugu.product_price_snapshot
WHERE product_id = $1
  AND currency = $2
  AND snapshot_date >= $3
  AND snapshot_date <= $4
ORDER BY snapshot_date ASC;
