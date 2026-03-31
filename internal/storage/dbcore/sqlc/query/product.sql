-- name: CreateProduct :exec
INSERT INTO gugu.product (
    id,
    market,
    external_product_id,
    original_url,
    title,
    main_image_url,
    current_price,
    currency,
    product_url,
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
);

-- name: UpdateProduct :execrows
UPDATE gugu.product
SET
    original_url = $2,
    title = $3,
    main_image_url = $4,
    current_price = $5,
    currency = $6,
    product_url = $7,
    promotion_link = $8,
    collection_source = $9,
    last_collected_at = $10,
    updated_at = $11
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
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product
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
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product
WHERE market = $1 AND external_product_id = $2;

-- name: ListProductsByMarket :many
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
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product
WHERE market = $1
ORDER BY created_at;

-- name: FindProductsByIDs :many
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
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product
WHERE id = ANY($1::text[]);
