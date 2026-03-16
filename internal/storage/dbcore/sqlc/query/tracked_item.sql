-- name: CreateTrackedItem :exec
INSERT INTO gugu.user_tracked_items (
    id,
    user_id,
    product_id,
    original_url,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: FindTrackedItemByUserIDAndProductID :one
SELECT
    id,
    user_id,
    product_id,
    original_url,
    deleted_at,
    created_at
FROM gugu.user_tracked_items
WHERE user_id = $1 AND product_id = $2 AND deleted_at IS NULL;

-- name: FindTrackedItemByIDAndUserID :one
SELECT
    id,
    user_id,
    product_id,
    original_url,
    deleted_at,
    created_at
FROM gugu.user_tracked_items
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: ListTrackedItemsByUserID :many
SELECT
    id,
    user_id,
    product_id,
    original_url,
    deleted_at,
    created_at
FROM gugu.user_tracked_items
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteTrackedItemByIDAndUserID :execrows
UPDATE gugu.user_tracked_items
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
