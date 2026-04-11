-- Flyway Migration
-- Version: V20260411_02
-- Description: backfill for tracked item view id, product/sku localization, product external alias

-- ============================================================================
-- 1) user_tracked_item.view_external_product_id backfill
--    우선순위:
--    (a) original_url 에서 /item/{id} 패턴 추출
--    (b) original_url 에서 /i/{id}.html 패턴 추출
--    (c) product.external_product_id fallback
-- ============================================================================
UPDATE gugu.user_tracked_item t
SET view_external_product_id = COALESCE(
    NULLIF(substring(t.original_url FROM '/item/([0-9]+)'), ''),
    NULLIF(substring(t.original_url FROM '/i/([0-9]+)\\.html'), ''),
    p.external_product_id,
    ''
)
FROM gugu.product p
WHERE t.product_id = p.id
  AND COALESCE(TRIM(t.view_external_product_id), '') = '';


-- ============================================================================
-- 2) product_localization backfill
--    2-1) product_variant 기반 (언어별 메타가 있는 경우)
--    2-2) 미존재 product는 product 본문으로 KO 기본 row 보강
-- ============================================================================
INSERT INTO gugu.product_localization (
    product_id,
    language,
    title,
    main_image_url,
    product_url,
    updated_at
)
SELECT
    x.product_id,
    x.language,
    x.title,
    x.main_image_url,
    x.product_url,
    NOW()
FROM (
    SELECT DISTINCT ON (pv.product_id, pv.language)
        pv.product_id,
        pv.language,
        COALESCE(pv.title, '') AS title,
        COALESCE(pv.main_image_url, '') AS main_image_url,
        COALESCE(pv.product_url, '') AS product_url
    FROM gugu.product_variant pv
    ORDER BY pv.product_id, pv.language, pv.updated_at DESC, pv.created_at DESC
) x
ON CONFLICT (product_id, language)
DO UPDATE SET
    title = EXCLUDED.title,
    main_image_url = EXCLUDED.main_image_url,
    product_url = EXCLUDED.product_url,
    updated_at = NOW();

INSERT INTO gugu.product_localization (
    product_id,
    language,
    title,
    main_image_url,
    product_url,
    updated_at
)
SELECT
    p.id,
    'KO',
    COALESCE(p.title, ''),
    COALESCE(p.main_image_url, ''),
    COALESCE(p.product_url, ''),
    NOW()
FROM gugu.product p
WHERE NOT EXISTS (
    SELECT 1
    FROM gugu.product_localization pl
    WHERE pl.product_id = p.id
)
ON CONFLICT (product_id, language)
DO UPDATE SET
    title = EXCLUDED.title,
    main_image_url = EXCLUDED.main_image_url,
    product_url = EXCLUDED.product_url,
    updated_at = NOW();


-- ============================================================================
-- 3) sku_localization backfill
--    현재 sku 메타를 KO 기본 row로 적재
-- ============================================================================
INSERT INTO gugu.sku_localization (
    sku_id,
    language,
    sku_name,
    color_name,
    size_name,
    sku_properties,
    image_url,
    updated_at
)
SELECT
    s.id,
    'KO',
    COALESCE(s.sku_name, ''),
    COALESCE(s.color, ''),
    COALESCE(s.size, ''),
    COALESCE(s.sku_properties, ''),
    COALESCE(s.image_url, ''),
    NOW()
FROM gugu.sku s
ON CONFLICT (sku_id, language)
DO UPDATE SET
    sku_name = EXCLUDED.sku_name,
    color_name = EXCLUDED.color_name,
    size_name = EXCLUDED.size_name,
    sku_properties = EXCLUDED.sku_properties,
    image_url = EXCLUDED.image_url,
    updated_at = NOW();


-- ============================================================================
-- 4) product_external_alias backfill (view id -> canonical product)
--    id는 deterministic hash 사용
-- ============================================================================
INSERT INTO gugu.product_external_alias (
    id,
    market,
    alias_external_product_id,
    product_id,
    alias_type,
    created_at,
    updated_at
)
SELECT
    md5(CONCAT(x.market, ':', x.view_external_product_id, ':VIEW')) AS id,
    x.market,
    x.view_external_product_id,
    x.product_id,
    'VIEW',
    NOW(),
    NOW()
FROM (
    SELECT DISTINCT ON (p.market, t.view_external_product_id)
        p.market,
        t.view_external_product_id,
        t.product_id,
        t.created_at
    FROM gugu.user_tracked_item t
    JOIN gugu.product p ON p.id = t.product_id
    WHERE COALESCE(TRIM(t.view_external_product_id), '') <> ''
    ORDER BY p.market, t.view_external_product_id, t.created_at DESC
) x
ON CONFLICT (market, alias_external_product_id)
DO UPDATE SET
    product_id = EXCLUDED.product_id,
    alias_type = EXCLUDED.alias_type,
    updated_at = NOW();


-- ============================================================================
-- 5) tracking_scope 기본값 보정
-- ============================================================================
UPDATE gugu.user_tracked_item
SET tracking_scope = 'PRODUCT_ALL_SKU'
WHERE COALESCE(TRIM(tracking_scope), '') = '';
