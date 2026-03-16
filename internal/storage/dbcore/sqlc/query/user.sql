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
