-- name: CreateTrackedItem :exec
INSERT INTO gugu.user_tracked_item (
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    currency,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: FindTrackedItemByUserIDAndProductID :one
SELECT
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    currency,
    deleted_at,
    created_at
FROM gugu.user_tracked_item
WHERE user_id = $1 AND product_id = $2 AND deleted_at IS NULL;

-- name: FindTrackedItemByIDAndUserID :one
SELECT
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    currency,
    deleted_at,
    created_at
FROM gugu.user_tracked_item
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: ListTrackedItemsByUserID :many
SELECT
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    currency,
    deleted_at,
    created_at
FROM gugu.user_tracked_item
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteTrackedItemByIDAndUserID :execrows
UPDATE gugu.user_tracked_item
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: UpdateTrackedItemSKU :execrows
UPDATE gugu.user_tracked_item
SET sku_id = $3
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
