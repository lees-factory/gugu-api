package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

type LoginSessionSQLRepository struct {
	db *sql.DB
}

func NewLoginSessionRepository(db *sql.DB) *LoginSessionSQLRepository {
	return &LoginSessionSQLRepository{db: db}
}

func (r *LoginSessionSQLRepository) Create(ctx context.Context, session domainauth.LoginSession) error {
	const query = `
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
)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.RefreshTokenHash,
		session.TokenFamilyID,
		session.ParentSessionID,
		session.UserAgent,
		session.ClientIP,
		session.DeviceName,
		session.ExpiresAt,
		session.LastSeenAt,
		session.RotatedAt,
		session.RevokedAt,
		session.ReuseDetectedAt,
		session.CreatedAt,
	)
	return err
}

func (r *LoginSessionSQLRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domainauth.LoginSession, error) {
	const query = `
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
WHERE refresh_token_hash = $1`

	var session domainauth.LoginSession
	var parentSessionID sql.NullString
	var rotatedAt sql.NullTime
	var revokedAt sql.NullTime
	var reuseDetectedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, refreshTokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.TokenFamilyID,
		&parentSessionID,
		&session.UserAgent,
		&session.ClientIP,
		&session.DeviceName,
		&session.ExpiresAt,
		&session.LastSeenAt,
		&rotatedAt,
		&revokedAt,
		&reuseDetectedAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if parentSessionID.Valid {
		session.ParentSessionID = &parentSessionID.String
	}
	if rotatedAt.Valid {
		session.RotatedAt = &rotatedAt.Time
	}
	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.Time
	}
	if reuseDetectedAt.Valid {
		session.ReuseDetectedAt = &reuseDetectedAt.Time
	}

	return &session, nil
}

func (r *LoginSessionSQLRepository) MarkRotated(ctx context.Context, sessionID string, rotatedAt time.Time) error {
	return r.execUpdate(ctx, `UPDATE gugu.user_login_sessions SET rotated_at = $2 WHERE id = $1`, sessionID, rotatedAt)
}

func (r *LoginSessionSQLRepository) Revoke(ctx context.Context, sessionID string, revokedAt time.Time) error {
	return r.execUpdate(ctx, `UPDATE gugu.user_login_sessions SET revoked_at = $2 WHERE id = $1`, sessionID, revokedAt)
}

func (r *LoginSessionSQLRepository) RevokeFamily(ctx context.Context, tokenFamilyID string, revokedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE gugu.user_login_sessions SET revoked_at = $2 WHERE token_family_id = $1`, tokenFamilyID, revokedAt)
	return err
}

func (r *LoginSessionSQLRepository) MarkReuseDetected(ctx context.Context, sessionID string, detectedAt time.Time) error {
	return r.execUpdate(ctx, `UPDATE gugu.user_login_sessions SET reuse_detected_at = $2 WHERE id = $1`, sessionID, detectedAt)
}

func (r *LoginSessionSQLRepository) UpdateLastSeen(ctx context.Context, sessionID string, lastSeenAt time.Time) error {
	return r.execUpdate(ctx, `UPDATE gugu.user_login_sessions SET last_seen_at = $2 WHERE id = $1`, sessionID, lastSeenAt)
}

func (r *LoginSessionSQLRepository) execUpdate(ctx context.Context, query string, key string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, query, key, at)
	return err
}
