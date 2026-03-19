-- name: CreateEmailVerification :exec
INSERT INTO gugu.email_verification (
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
FROM gugu.email_verification
WHERE code = $1;

-- name: MarkEmailVerificationUsed :execrows
UPDATE gugu.email_verification
SET used_at = $2
WHERE code = $1;
