-- name: CreateUser :exec
INSERT INTO gugu.app_users (
    id,
    email,
    display_name,
    password_hash,
    auth_source,
    email_verified,
    email_verified_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: FindUserByEmail :one
SELECT
    id,
    email,
    display_name,
    password_hash,
    auth_source,
    email_verified,
    email_verified_at,
    created_at
FROM gugu.app_users
WHERE email = $1;

-- name: FindUserByID :one
SELECT
    id,
    email,
    display_name,
    password_hash,
    auth_source,
    email_verified,
    email_verified_at,
    created_at
FROM gugu.app_users
WHERE id = $1;

-- name: MarkUserEmailVerified :execrows
UPDATE gugu.app_users
SET
    email_verified = TRUE,
    email_verified_at = $2
WHERE id = $1;

-- name: CreateEmailVerification :exec
INSERT INTO gugu.email_verifications (
    code,
    user_id,
    email,
    expires_at,
    used_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: FindEmailVerificationByCode :one
SELECT
    code,
    user_id,
    email,
    expires_at,
    used_at,
    created_at
FROM gugu.email_verifications
WHERE code = $1;

-- name: MarkEmailVerificationUsed :execrows
UPDATE gugu.email_verifications
SET used_at = $2
WHERE code = $1;

-- name: CreateOAuthIdentity :exec
INSERT INTO gugu.oauth_identities (
    id,
    user_id,
    provider,
    subject,
    email,
    created_at,
    last_login_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: FindOAuthIdentity :one
SELECT
    id,
    user_id,
    provider,
    subject,
    email,
    created_at,
    last_login_at
FROM gugu.oauth_identities
WHERE provider = $1 AND subject = $2;

-- name: UpdateOAuthLastLogin :execrows
UPDATE gugu.oauth_identities
SET last_login_at = $3
WHERE provider = $1 AND subject = $2;

-- name: CreateUserLoginSession :exec
INSERT INTO gugu.user_login_sessions (
    id,
    user_id,
    refresh_token_hash,
    token_family_id,
    parent_session_id,
    user_agent,
    client_ip,
    device_name,
    expires_at,
    last_seen_at,
    rotated_at,
    revoked_at,
    reuse_detected_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
);

-- name: FindUserLoginSessionByRefreshTokenHash :one
SELECT
    id,
    user_id,
    refresh_token_hash,
    token_family_id,
    parent_session_id,
    user_agent,
    client_ip,
    device_name,
    expires_at,
    last_seen_at,
    rotated_at,
    revoked_at,
    reuse_detected_at,
    created_at
FROM gugu.user_login_sessions
WHERE refresh_token_hash = $1;

-- name: MarkUserLoginSessionRotated :execrows
UPDATE gugu.user_login_sessions
SET rotated_at = $2
WHERE id = $1;

-- name: RevokeUserLoginSession :execrows
UPDATE gugu.user_login_sessions
SET revoked_at = $2
WHERE id = $1;

-- name: RevokeUserLoginSessionFamily :execrows
UPDATE gugu.user_login_sessions
SET revoked_at = $2
WHERE token_family_id = $1;

-- name: MarkUserLoginSessionReuseDetected :execrows
UPDATE gugu.user_login_sessions
SET reuse_detected_at = $2
WHERE id = $1;

-- name: UpdateUserLoginSessionLastSeen :execrows
UPDATE gugu.user_login_sessions
SET last_seen_at = $2
WHERE id = $1;

-- name: UpsertAliExpressSellerToken :exec
INSERT INTO gugu.aliexpress_seller_tokens (
    id,
    user_id,
    seller_id,
    havana_id,
    app_user_id,
    user_nick,
    account,
    account_platform,
    locale,
    sp,
    access_token,
    refresh_token,
    access_token_expires_at,
    refresh_token_expires_at,
    last_refreshed_at,
    authorized_at,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
)
ON CONFLICT (seller_id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    havana_id = EXCLUDED.havana_id,
    app_user_id = EXCLUDED.app_user_id,
    user_nick = EXCLUDED.user_nick,
    account = EXCLUDED.account,
    account_platform = EXCLUDED.account_platform,
    locale = EXCLUDED.locale,
    sp = EXCLUDED.sp,
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    access_token_expires_at = EXCLUDED.access_token_expires_at,
    refresh_token_expires_at = EXCLUDED.refresh_token_expires_at,
    last_refreshed_at = EXCLUDED.last_refreshed_at,
    authorized_at = EXCLUDED.authorized_at,
    updated_at = EXCLUDED.updated_at;

-- name: FindAliExpressSellerTokenByUserID :one
SELECT
    id,
    user_id,
    seller_id,
    havana_id,
    app_user_id,
    user_nick,
    account,
    account_platform,
    locale,
    sp,
    access_token,
    refresh_token,
    access_token_expires_at,
    refresh_token_expires_at,
    last_refreshed_at,
    authorized_at,
    created_at,
    updated_at
FROM gugu.aliexpress_seller_tokens
WHERE user_id = $1;

-- name: FindAliExpressSellerTokenBySellerID :one
SELECT
    id,
    user_id,
    seller_id,
    havana_id,
    app_user_id,
    user_nick,
    account,
    account_platform,
    locale,
    sp,
    access_token,
    refresh_token,
    access_token_expires_at,
    refresh_token_expires_at,
    last_refreshed_at,
    authorized_at,
    created_at,
    updated_at
FROM gugu.aliexpress_seller_tokens
WHERE seller_id = $1;

-- name: ListAliExpressSellerTokensExpiringBefore :many
SELECT
    id,
    user_id,
    seller_id,
    havana_id,
    app_user_id,
    user_nick,
    account,
    account_platform,
    locale,
    sp,
    access_token,
    refresh_token,
    access_token_expires_at,
    refresh_token_expires_at,
    last_refreshed_at,
    authorized_at,
    created_at,
    updated_at
FROM gugu.aliexpress_seller_tokens
WHERE access_token_expires_at <= $1
ORDER BY access_token_expires_at ASC;
