-- aliexpress_seller_token: app_type 컬럼 추가 + UNIQUE 제약 변경
ALTER TABLE gugu.aliexpress_seller_token
    ADD COLUMN IF NOT EXISTS app_type TEXT NOT NULL DEFAULT 'AFFILIATE';

ALTER TABLE gugu.aliexpress_seller_token
    DROP CONSTRAINT IF EXISTS aliexpress_seller_token_seller_id_key;

ALTER TABLE gugu.aliexpress_seller_token
    ADD CONSTRAINT uq_aliexpress_seller_token_app_type_seller_id UNIQUE (app_type, seller_id);
