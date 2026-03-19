-- name: CreateUserLoginSession :exec
INSERT INTO gugu.user_login_session (
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
FROM gugu.user_login_session
WHERE refresh_token_hash = $1;

-- name: MarkUserLoginSessionRotated :execrows
UPDATE gugu.user_login_session
SET rotated_at = $2
WHERE id = $1;

-- name: RevokeUserLoginSession :execrows
UPDATE gugu.user_login_session
SET revoked_at = $2
WHERE id = $1;

-- name: RevokeUserLoginSessionFamily :execrows
UPDATE gugu.user_login_session
SET revoked_at = $2
WHERE token_family_id = $1;

-- name: MarkUserLoginSessionReuseDetected :execrows
UPDATE gugu.user_login_session
SET reuse_detected_at = $2
WHERE id = $1;

-- name: UpdateUserLoginSessionLastSeen :execrows
UPDATE gugu.user_login_session
SET last_seen_at = $2
WHERE id = $1;
