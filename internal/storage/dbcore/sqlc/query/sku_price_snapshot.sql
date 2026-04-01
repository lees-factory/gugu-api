-- name: UpsertSKUPriceSnapshot :exec
INSERT INTO gugu.sku_price_snapshot (
    sku_id,
    snapshot_date,
    price,
    original_price,
    currency
) VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (sku_id, currency, snapshot_date) DO UPDATE SET
    price = EXCLUDED.price,
    original_price = EXCLUDED.original_price;

-- name: ListSKUPriceSnapshotsByDateRange :many
SELECT
    sku_id,
    snapshot_date,
    price,
    original_price,
    currency
FROM gugu.sku_price_snapshot
WHERE sku_id = $1
  AND currency = $2
  AND snapshot_date >= $3
  AND snapshot_date <= $4
ORDER BY snapshot_date ASC;
