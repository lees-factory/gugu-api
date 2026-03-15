CREATE SCHEMA IF NOT EXISTS gugu;

CREATE TABLE IF NOT EXISTS gugu.app_users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL DEFAULT '',
    auth_source TEXT NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gugu.email_verifications (
    code TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_users(id),
    email TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gugu.oauth_identities (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_users(id),
    provider TEXT NOT NULL,
    subject TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ NOT NULL,
    UNIQUE (provider, subject)
);

CREATE TABLE IF NOT EXISTS gugu.user_login_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_users(id),
    refresh_token_hash TEXT NOT NULL UNIQUE,
    token_family_id TEXT NOT NULL,
    parent_session_id TEXT REFERENCES gugu.user_login_sessions(id),
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

CREATE TABLE IF NOT EXISTS gugu.aliexpress_seller_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES gugu.app_users(id),
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

CREATE INDEX IF NOT EXISTS idx_app_users_email ON gugu.app_users(email);
CREATE INDEX IF NOT EXISTS idx_email_verifications_user_id ON gugu.email_verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_identities_user_id ON gugu.oauth_identities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_login_sessions_user_id ON gugu.user_login_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_login_sessions_family_id ON gugu.user_login_sessions(token_family_id);
CREATE INDEX IF NOT EXISTS idx_aliexpress_seller_tokens_user_id ON gugu.aliexpress_seller_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_aliexpress_seller_tokens_access_token_expires_at ON gugu.aliexpress_seller_tokens(access_token_expires_at);
CREATE INDEX IF NOT EXISTS idx_aliexpress_seller_tokens_refresh_token_expires_at ON gugu.aliexpress_seller_tokens(refresh_token_expires_at);
