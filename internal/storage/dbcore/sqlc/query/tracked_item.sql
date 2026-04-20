-- name: CreateTrackedItem :exec
INSERT INTO gugu.user_tracked_item (
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    view_external_product_id,
    preferred_language,
    tracking_scope,
    currency,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: FindTrackedItemByUserIDAndProductID :one
SELECT
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    view_external_product_id,
    preferred_language,
    tracking_scope,
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
    view_external_product_id,
    preferred_language,
    tracking_scope,
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
    view_external_product_id,
    preferred_language,
    tracking_scope,
    currency,
    deleted_at,
    created_at
FROM gugu.user_tracked_item
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListTrackedItemsByUserIDWithCursor :many
SELECT
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    view_external_product_id,
    preferred_language,
    tracking_scope,
    currency,
    deleted_at,
    created_at
FROM gugu.user_tracked_item
WHERE user_id = $1
  AND deleted_at IS NULL
  AND (created_at < $2 OR (created_at = $2 AND id < $3))
ORDER BY created_at DESC, id DESC
LIMIT $4;

-- name: ListTrackedItemsByUserIDFirstPage :many
SELECT
    id,
    user_id,
    product_id,
    sku_id,
    original_url,
    view_external_product_id,
    preferred_language,
    tracking_scope,
    currency,
    deleted_at,
    created_at
FROM gugu.user_tracked_item
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT $2;

-- name: DeleteTrackedItemByIDAndUserID :execrows
UPDATE gugu.user_tracked_item
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: UpdateTrackedItemSKU :execrows
UPDATE gugu.user_tracked_item
SET sku_id = $3
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: UpdateTrackedItemPreferredLanguage :execrows
UPDATE gugu.user_tracked_item
SET preferred_language = $3
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: UpdateTrackedItemTrackingScope :execrows
UPDATE gugu.user_tracked_item
SET tracking_scope = $3
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
