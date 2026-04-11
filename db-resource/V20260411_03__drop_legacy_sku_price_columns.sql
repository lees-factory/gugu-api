-- Flyway Migration
-- Version: V20260411_03
-- Description: drop legacy sku direct price columns after cutover validation
--
-- Pre-check (must be satisfied before applying):
-- 1) API read path uses snapshot/history as source of truth for price fields.
-- 2) Batch recorder/publisher no longer depends on gugu.sku.price/original_price/currency.
-- 3) Phase 2/4 validation gate is GO in production-like environment.

ALTER TABLE gugu.sku
    DROP COLUMN IF EXISTS price;

ALTER TABLE gugu.sku
    DROP COLUMN IF EXISTS original_price;

ALTER TABLE gugu.sku
    DROP COLUMN IF EXISTS currency;
