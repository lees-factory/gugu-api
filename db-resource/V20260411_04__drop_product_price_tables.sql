-- Flyway Migration
-- Version: V20260411_04
-- Description: drop legacy product-level price tables after sku-only cutover
--
-- Pre-check (must be satisfied before applying):
-- 1) client-api no longer reads product_price_history/product_price_snapshot.
-- 2) product trend API is sku_id 기반으로만 조회된다.
-- 3) monitoring confirms no runtime SQL hits to product price tables.

DROP TABLE IF EXISTS gugu.product_price_history;
DROP TABLE IF EXISTS gugu.product_price_snapshot;
