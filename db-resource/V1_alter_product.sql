-- product: promotion_link 컬럼 추가
ALTER TABLE gugu.product
    ADD COLUMN IF NOT EXISTS promotion_link TEXT NOT NULL DEFAULT '';

-- product: current_price 컬럼 제거
ALTER TABLE gugu.product
    DROP COLUMN IF EXISTS current_price;

-- product: currency 컬럼 제거
ALTER TABLE gugu.product
    DROP COLUMN IF EXISTS currency;
