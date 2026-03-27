-- name: UpsertProductPriceSnapshot :exec
INSERT INTO gugu.product_price_snapshot (
    product_id,
    snapshot_date,
    price,
    currency
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (product_id, snapshot_date) DO UPDATE SET
    price = EXCLUDED.price,
    currency = EXCLUDED.currency;

-- name: ListProductPriceSnapshotsByDateRange :many
SELECT
    product_id,
    snapshot_date,
    price,
    currency
FROM gugu.product_price_snapshot
WHERE product_id = $1
  AND snapshot_date >= $2
  AND snapshot_date <= $3
ORDER BY snapshot_date ASC;
