-- name: CreateUser :exec
INSERT INTO gugu.app_user (
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
FROM gugu.app_user
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
FROM gugu.app_user
WHERE id = $1;

-- name: MarkUserEmailVerified :execrows
UPDATE gugu.app_user
SET
    email_verified = TRUE,
    email_verified_at = $2
WHERE id = $1;
