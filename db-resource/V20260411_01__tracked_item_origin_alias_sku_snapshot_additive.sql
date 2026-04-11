-- Flyway Migration
-- Version: V20260411_01
-- Description: tracked item origin alias / localization / selected sku / snapshot run additive schema

-- 1) product alias: view id -> canonical product
CREATE TABLE IF NOT EXISTS gugu.product_external_alias (
    id TEXT PRIMARY KEY,
    market TEXT NOT NULL,
    alias_external_product_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    alias_type TEXT NOT NULL DEFAULT 'VIEW', -- VIEW | ORIGIN_REDIRECT
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_external_alias_product
        FOREIGN KEY (product_id) REFERENCES gugu.product(id),
    CONSTRAINT uq_product_external_alias_market_alias
        UNIQUE (market, alias_external_product_id)
);

CREATE INDEX IF NOT EXISTS idx_product_external_alias_product_id
    ON gugu.product_external_alias(product_id);


-- 2) product localization: language specific metadata
CREATE TABLE IF NOT EXISTS gugu.product_localization (
    product_id TEXT NOT NULL,
    language TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    main_image_url TEXT NOT NULL DEFAULT '',
    product_url TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (product_id, language),
    CONSTRAINT fk_product_localization_product
        FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);


-- 3) sku localization: language specific sku metadata
CREATE TABLE IF NOT EXISTS gugu.sku_localization (
    sku_id TEXT NOT NULL,
    language TEXT NOT NULL,
    sku_name TEXT NOT NULL DEFAULT '',
    color_name TEXT NOT NULL DEFAULT '',
    size_name TEXT NOT NULL DEFAULT '',
    sku_properties TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (sku_id, language),
    CONSTRAINT fk_sku_localization_sku
        FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);


-- 4) tracked item extension: canonical + view + display preference + scope
ALTER TABLE gugu.user_tracked_item
    ADD COLUMN IF NOT EXISTS view_external_product_id TEXT NOT NULL DEFAULT '';

ALTER TABLE gugu.user_tracked_item
    ADD COLUMN IF NOT EXISTS preferred_language TEXT NOT NULL DEFAULT 'KO';

ALTER TABLE gugu.user_tracked_item
    ADD COLUMN IF NOT EXISTS tracking_scope TEXT NOT NULL DEFAULT 'PRODUCT_ALL_SKU';
-- PRODUCT_ALL_SKU | SELECTED_SKU_ONLY

CREATE INDEX IF NOT EXISTS idx_user_tracked_item_view_external_product_id
    ON gugu.user_tracked_item(view_external_product_id);


-- 5) selected SKU tracking table
CREATE TABLE IF NOT EXISTS gugu.user_tracked_item_watch_sku (
    tracked_item_id TEXT NOT NULL,
    sku_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tracked_item_id, sku_id),
    CONSTRAINT fk_user_tracked_item_watch_sku_tracked_item
        FOREIGN KEY (tracked_item_id) REFERENCES gugu.user_tracked_item(id),
    CONSTRAINT fk_user_tracked_item_watch_sku_sku
        FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);

CREATE INDEX IF NOT EXISTS idx_tracked_item_watch_sku_sku_id
    ON gugu.user_tracked_item_watch_sku(sku_id);


-- 6) snapshot ingest run + staging table for atomic publish
CREATE TABLE IF NOT EXISTS gugu.sku_snapshot_ingest_run (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL,
    currency TEXT NOT NULL,
    snapshot_date DATE NOT NULL,
    expected_sku_count INTEGER NOT NULL DEFAULT 0,
    collected_sku_count INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL, -- PENDING|RUNNING|COMPLETED|PARTIAL|FAILED
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    error_message TEXT NOT NULL DEFAULT '',
    CONSTRAINT fk_sku_snapshot_ingest_run_product
        FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);

CREATE INDEX IF NOT EXISTS idx_sku_snapshot_ingest_run_product_currency_date
    ON gugu.sku_snapshot_ingest_run(product_id, currency, snapshot_date DESC);

CREATE TABLE IF NOT EXISTS gugu.sku_price_snapshot_staging (
    run_id TEXT NOT NULL,
    sku_id TEXT NOT NULL,
    snapshot_date DATE NOT NULL,
    price TEXT NOT NULL DEFAULT '',
    original_price TEXT NOT NULL DEFAULT '',
    currency TEXT NOT NULL,
    PRIMARY KEY (run_id, sku_id, currency),
    CONSTRAINT fk_sku_price_snapshot_staging_run
        FOREIGN KEY (run_id) REFERENCES gugu.sku_snapshot_ingest_run(id),
    CONSTRAINT fk_sku_price_snapshot_staging_sku
        FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);


-- 7) deprecate sku direct price columns (cutover 이후 drop 예정)
COMMENT ON COLUMN gugu.sku.price IS
'@deprecated: 가격 정본은 sku_price_snapshot/sku_price_history 사용';
COMMENT ON COLUMN gugu.sku.original_price IS
'@deprecated: 가격 정본은 sku_price_snapshot/sku_price_history 사용';
COMMENT ON COLUMN gugu.sku.currency IS
'@deprecated: 가격 정본은 sku_price_snapshot/sku_price_history 사용';
