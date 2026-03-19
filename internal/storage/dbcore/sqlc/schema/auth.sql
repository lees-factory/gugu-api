CREATE SCHEMA IF NOT EXISTS gugu;

CREATE TABLE IF NOT EXISTS gugu.app_user (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL DEFAULT '',
    auth_source TEXT NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gugu.email_verification (
    code TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_user(id),
    email TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gugu.oauth_identity (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_user(id),
    provider TEXT NOT NULL,
    subject TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ NOT NULL,
    UNIQUE (provider, subject)
);

CREATE TABLE IF NOT EXISTS gugu.user_login_session (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_user(id),
    refresh_token_hash TEXT NOT NULL UNIQUE,
    token_family_id TEXT NOT NULL,
    parent_session_id TEXT REFERENCES gugu.user_login_session(id),
    user_agent TEXT NOT NULL DEFAULT '',
    client_ip TEXT NOT NULL DEFAULT '',
    device_name TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    rotated_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    reuse_detected_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gugu.aliexpress_seller_token (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_user(id),
    seller_id TEXT NOT NULL UNIQUE,
    havana_id TEXT NOT NULL DEFAULT '',
    app_user_id TEXT NOT NULL DEFAULT '',
    user_nick TEXT NOT NULL DEFAULT '',
    account TEXT NOT NULL DEFAULT '',
    account_platform TEXT NOT NULL DEFAULT '',
    locale TEXT NOT NULL DEFAULT '',
    sp TEXT NOT NULL DEFAULT '',
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    access_token_expires_at TIMESTAMPTZ NOT NULL,
    refresh_token_expires_at TIMESTAMPTZ,
    last_refreshed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    authorized_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gugu.product (
    id TEXT PRIMARY KEY,
    market TEXT NOT NULL,
    external_product_id TEXT NOT NULL,
    original_url TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    main_image_url TEXT NOT NULL DEFAULT '',
    current_price TEXT NOT NULL DEFAULT '',
    currency TEXT NOT NULL DEFAULT '',
    product_url TEXT NOT NULL DEFAULT '',
    collection_source TEXT NOT NULL DEFAULT '',
    last_collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (market, external_product_id)
);

CREATE TABLE IF NOT EXISTS gugu.user_tracked_item (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_user(id),
    product_id TEXT NOT NULL REFERENCES gugu.product(id),
    sku_id TEXT,
    original_url TEXT NOT NULL DEFAULT '',
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gugu.product_price_history (
    product_id TEXT NOT NULL REFERENCES gugu.product(id),
    recorded_at TIMESTAMPTZ NOT NULL,
    price TEXT NOT NULL DEFAULT '',
    currency TEXT NOT NULL DEFAULT '',
    change_value TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (product_id, recorded_at)
);

CREATE INDEX IF NOT EXISTS idx_app_user_email ON gugu.app_user(email);
CREATE INDEX IF NOT EXISTS idx_email_verification_user_id ON gugu.email_verification(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_identity_user_id ON gugu.oauth_identity(user_id);
CREATE INDEX IF NOT EXISTS idx_user_login_session_user_id ON gugu.user_login_session(user_id);
CREATE INDEX IF NOT EXISTS idx_user_login_session_family_id ON gugu.user_login_session(token_family_id);
CREATE INDEX IF NOT EXISTS idx_aliexpress_seller_token_user_id ON gugu.aliexpress_seller_token(user_id);
CREATE INDEX IF NOT EXISTS idx_aliexpress_seller_token_access_token_expires_at ON gugu.aliexpress_seller_token(access_token_expires_at);
CREATE INDEX IF NOT EXISTS idx_aliexpress_seller_token_refresh_token_expires_at ON gugu.aliexpress_seller_token(refresh_token_expires_at);
CREATE INDEX IF NOT EXISTS idx_product_market_external_product_id ON gugu.product(market, external_product_id);
CREATE INDEX IF NOT EXISTS idx_user_tracked_item_user_id ON gugu.user_tracked_item(user_id);
CREATE INDEX IF NOT EXISTS idx_user_tracked_item_product_id ON gugu.user_tracked_item(product_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_user_tracked_item_user_product_active
    ON gugu.user_tracked_item(user_id, product_id)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_product_price_history_product_id_recorded_at ON gugu.product_price_history(product_id, recorded_at DESC);

CREATE TABLE IF NOT EXISTS gugu.product_sku (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL REFERENCES gugu.product(id),
    external_sku_id TEXT NOT NULL DEFAULT '',
    sku_name TEXT NOT NULL DEFAULT '',
    color TEXT NOT NULL DEFAULT '',
    size TEXT NOT NULL DEFAULT '',
    price TEXT NOT NULL DEFAULT '',
    original_price TEXT NOT NULL DEFAULT '',
    currency TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    sku_properties TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, external_sku_id)
);

CREATE INDEX IF NOT EXISTS idx_product_sku_product_id ON gugu.product_sku(product_id);
