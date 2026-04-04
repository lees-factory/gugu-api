-- name: CreateProduct :exec
INSERT INTO gugu.product (
    id,
    market,
    external_product_id,
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
    external_product_id,
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
    external_product_id,
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
WHERE market = $1 AND external_product_id = $2;

-- name: ListProductsByMarket :many
SELECT
    id,
    market,
    external_product_id,
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
    external_product_id,
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
    external_product_id,
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

-- name: FindProductVariantByProductIDLanguageCurrency :one
SELECT
    product_id,
    language,
    currency,
    title,
    main_image_url,
    product_url,
    current_price,
    last_collected_at,
    created_at,
    updated_at
FROM gugu.product_variant
WHERE product_id = $1
  AND language = $2
  AND currency = $3;

-- name: FindProductVariantsByLookupKeys :many
WITH args AS (
    SELECT
        $1::text[] AS product_ids,
        $2::text[] AS languages,
        $3::text[] AS currencies
),
lookup AS (
    SELECT
        product_ids[idx] AS product_id,
        languages[idx] AS language,
        currencies[idx] AS currency
    FROM args, generate_subscripts(product_ids, 1) AS idx
)
SELECT
    pv.product_id,
    pv.language,
    pv.currency,
    pv.title,
    pv.main_image_url,
    pv.product_url,
    pv.current_price,
    pv.last_collected_at,
    pv.created_at,
    pv.updated_at
FROM gugu.product_variant pv
JOIN lookup
  ON pv.product_id = lookup.product_id
 AND pv.language = lookup.language
 AND pv.currency = lookup.currency;

-- name: UpsertProductVariant :exec
INSERT INTO gugu.product_variant (
    product_id,
    language,
    currency,
    title,
    main_image_url,
    product_url,
    current_price,
    last_collected_at,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (product_id, language, currency)
DO UPDATE SET
    title = EXCLUDED.title,
    main_image_url = EXCLUDED.main_image_url,
    product_url = EXCLUDED.product_url,
    current_price = EXCLUDED.current_price,
    last_collected_at = EXCLUDED.last_collected_at,
    updated_at = EXCLUDED.updated_at;
