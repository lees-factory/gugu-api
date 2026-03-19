-- name: CreateProductSKU :exec
INSERT INTO gugu.product_sku (
    id,
    product_id,
    external_sku_id,
    sku_name,
    color,
    size,
    price,
    original_price,
    currency,
    image_url,
    sku_properties,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: FindProductSKUByID :one
SELECT
    id,
    product_id,
    external_sku_id,
    sku_name,
    color,
    size,
    price,
    original_price,
    currency,
    image_url,
    sku_properties,
    created_at,
    updated_at
FROM gugu.product_sku
WHERE id = $1;

-- name: FindProductSKUsByProductID :many
SELECT
    id,
    product_id,
    external_sku_id,
    sku_name,
    color,
    size,
    price,
    original_price,
    currency,
    image_url,
    sku_properties,
    created_at,
    updated_at
FROM gugu.product_sku
WHERE product_id = $1
ORDER BY created_at;
