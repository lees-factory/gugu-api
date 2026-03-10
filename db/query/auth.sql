-- name: CreateUser :exec
INSERT INTO app_users (
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
FROM app_users
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
FROM app_users
WHERE id = $1;

-- name: MarkUserEmailVerified :exec
UPDATE app_users
SET
    email_verified = TRUE,
    email_verified_at = $2
WHERE id = $1;

-- name: CreateEmailVerification :exec
INSERT INTO email_verifications (
    token,
    user_id,
    email,
    expires_at,
    used_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: FindEmailVerificationByToken :one
SELECT
    token,
    user_id,
    email,
    expires_at,
    used_at,
    created_at
FROM email_verifications
WHERE token = $1;

-- name: MarkEmailVerificationUsed :exec
UPDATE email_verifications
SET used_at = $2
WHERE token = $1;

-- name: CreateOAuthIdentity :exec
INSERT INTO oauth_identities (
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
FROM oauth_identities
WHERE provider = $1 AND subject = $2;

-- name: UpdateOAuthLastLogin :exec
UPDATE oauth_identities
SET last_login_at = $3
WHERE provider = $1 AND subject = $2;

-- name: CreateSession :exec
INSERT INTO auth_sessions (
    token,
    user_id,
    expires_at,
    created_at
) VALUES (
    $1, $2, $3, $4
);
