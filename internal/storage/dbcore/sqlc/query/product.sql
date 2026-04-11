-- name: CreateProduct :exec
INSERT INTO gugu.product (
    id,
    market,
    origin_product_id,
    original_url,
    title,
    main_image_url,
    product_url,
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: UpdateProduct :execrows
UPDATE gugu.product
SET
    original_url = $2,
    title = $3,
    main_image_url = $4,
    product_url = $5,
    promotion_link = $6,
    collection_source = $7,
    last_collected_at = $8,
    updated_at = $9
WHERE id = $1;

-- name: FindProductByID :one
SELECT
    id,
    market,
    origin_product_id AS external_product_id,
    original_url,
    title,
    main_image_url,
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
    origin_product_id AS external_product_id,
    original_url,
    title,
    main_image_url,
    product_url,
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product
WHERE market = $1 AND origin_product_id = $2;

-- name: ListProductsByMarket :many
SELECT
    id,
    market,
    origin_product_id AS external_product_id,
    original_url,
    title,
    main_image_url,
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
    origin_product_id AS external_product_id,
    original_url,
    title,
    main_image_url,
    product_url,
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product
WHERE id = ANY($1::text[]);

-- name: ListProductsByCollectionSource :many
SELECT
    id,
    market,
    origin_product_id AS external_product_id,
    original_url,
    title,
    main_image_url,
    product_url,
    promotion_link,
    collection_source,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product
WHERE collection_source = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
