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
