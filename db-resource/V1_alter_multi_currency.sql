-- product_price_history: PK에 currency 추가
ALTER TABLE gugu.product_price_history DROP CONSTRAINT IF EXISTS product_price_history_pkey;
ALTER TABLE gugu.product_price_history DROP CONSTRAINT IF EXISTS gugu_product_price_history_pkey;
ALTER TABLE gugu.product_price_history DROP CONSTRAINT IF EXISTS "product_price_history_pkey1";
-- fallback: 이름 모를 경우 아래 쿼리로 확인
-- SELECT constraint_name FROM information_schema.table_constraints WHERE table_schema='gugu' AND table_name='product_price_history' AND constraint_type='PRIMARY KEY';
DO $$
DECLARE pk_name TEXT;
BEGIN
  SELECT constraint_name INTO pk_name
  FROM information_schema.table_constraints
  WHERE table_schema = 'gugu' AND table_name = 'product_price_history' AND constraint_type = 'PRIMARY KEY';
  IF pk_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE gugu.product_price_history DROP CONSTRAINT %I', pk_name);
  END IF;
END $$;
ALTER TABLE gugu.product_price_history ADD PRIMARY KEY (product_id, currency, recorded_at);

-- sku_price_history: PK에 currency 추가
DO $$
DECLARE pk_name TEXT;
BEGIN
  SELECT constraint_name INTO pk_name
  FROM information_schema.table_constraints
  WHERE table_schema = 'gugu' AND table_name = 'sku_price_history' AND constraint_type = 'PRIMARY KEY';
  IF pk_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE gugu.sku_price_history DROP CONSTRAINT %I', pk_name);
  END IF;
END $$;
ALTER TABLE gugu.sku_price_history ADD PRIMARY KEY (sku_id, currency, recorded_at);

-- product_price_snapshot: PK에 currency 추가
DO $$
DECLARE pk_name TEXT;
BEGIN
  SELECT constraint_name INTO pk_name
  FROM information_schema.table_constraints
  WHERE table_schema = 'gugu' AND table_name = 'product_price_snapshot' AND constraint_type = 'PRIMARY KEY';
  IF pk_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE gugu.product_price_snapshot DROP CONSTRAINT %I', pk_name);
  END IF;
END $$;
ALTER TABLE gugu.product_price_snapshot ADD PRIMARY KEY (product_id, currency, snapshot_date);
DROP INDEX IF EXISTS gugu.idx_product_price_snapshot_product_id;
CREATE INDEX IF NOT EXISTS idx_product_price_snapshot_lookup
    ON gugu.product_price_snapshot(product_id, currency, snapshot_date DESC);

-- sku_price_snapshot: PK에 currency 추가
DO $$
DECLARE pk_name TEXT;
BEGIN
  SELECT constraint_name INTO pk_name
  FROM information_schema.table_constraints
  WHERE table_schema = 'gugu' AND table_name = 'sku_price_snapshot' AND constraint_type = 'PRIMARY KEY';
  IF pk_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE gugu.sku_price_snapshot DROP CONSTRAINT %I', pk_name);
  END IF;
END $$;
ALTER TABLE gugu.sku_price_snapshot ADD PRIMARY KEY (sku_id, currency, snapshot_date);
DROP INDEX IF EXISTS gugu.idx_sku_price_snapshot_sku_id;
CREATE INDEX IF NOT EXISTS idx_sku_price_snapshot_lookup
    ON gugu.sku_price_snapshot(sku_id, currency, snapshot_date DESC);

-- user_tracked_item: currency 컬럼 추가
ALTER TABLE gugu.user_tracked_item ADD COLUMN IF NOT EXISTS currency TEXT NOT NULL DEFAULT 'KRW';
