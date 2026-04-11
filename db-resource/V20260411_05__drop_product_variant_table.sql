-- Flyway Migration
-- Version: V20260411_05
-- Description: drop legacy product_variant table after localization cutover
--
-- Pre-check (must be satisfied before applying):
-- 1) Variant/localized metadata read path uses gugu.product_localization.
-- 2) No runtime SQL hits to gugu.product_variant.
-- 3) Backfill to gugu.product_localization is complete.

DROP TABLE IF EXISTS gugu.product_variant;
