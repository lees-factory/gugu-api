-- name: CreateProduct :exec
INSERT INTO gugu.products (
    id,
    market,
    external_product_id,
    original_url,
    title,
    main_image_url,
    current_price,
    currency,
    product_url,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: UpdateProduct :execrows
UPDATE gugu.products
SET
    original_url = $2,
    title = $3,
    main_image_url = $4,
    current_price = $5,
    currency = $6,
    product_url = $7,
    collection_source = $8,
    last_collected_at = $9,
    updated_at = $10
WHERE id = $1;

-- name: FindProductByID :one
SELECT
    id,
    market,
    external_product_id,
    original_url,
    title,
    main_image_url,
    current_price,
    currency,
    product_url,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.products
WHERE id = $1;

-- name: FindProductByMarketAndExternalProductID :one
SELECT
    id,
    market,
    external_product_id,
    original_url,
    title,
    main_image_url,
    current_price,
    currency,
    product_url,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.products
WHERE market = $1 AND external_product_id = $2;
