CREATE TABLE IF NOT EXISTS gugu.product_price_snapshot (
    product_id TEXT NOT NULL,
    snapshot_date DATE NOT NULL,
    price TEXT NOT NULL DEFAULT '',
    currency TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (product_id, snapshot_date),
    CONSTRAINT fk_product_price_snapshot_product
        FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);

CREATE INDEX IF NOT EXISTS idx_product_price_snapshot_product_id
    ON gugu.product_price_snapshot(product_id, snapshot_date DESC);

CREATE TABLE IF NOT EXISTS gugu.sku_price_snapshot (
    sku_id TEXT NOT NULL,
    snapshot_date DATE NOT NULL,
    price TEXT NOT NULL DEFAULT '',
    original_price TEXT NOT NULL DEFAULT '',
    currency TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (sku_id, snapshot_date),
    CONSTRAINT fk_sku_price_snapshot_sku
        FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);

CREATE INDEX IF NOT EXISTS idx_sku_price_snapshot_sku_id
    ON gugu.sku_price_snapshot(sku_id, snapshot_date DESC);
