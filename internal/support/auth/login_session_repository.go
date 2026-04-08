package auth

import (
	"context"
	"time"
)

type LoginSessionRepository interface {
	Create(ctx context.Context, session LoginSession) error
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*LoginSession, error)
	ListActiveByUserID(ctx context.Context, userID string, now time.Time) ([]LoginSession, error)
	CountActiveByUserID(ctx context.Context, userID string, now time.Time) (int, error)
	MarkRotated(ctx context.Context, sessionID string, rotatedAt time.Time) error
	Revoke(ctx context.Context, sessionID string, revokedAt time.Time) error
	RevokeByUserSessionID(ctx context.Context, userID string, sessionID string, revokedAt time.Time) error
	RevokeOldestActiveByUserID(ctx context.Context, userID string, now time.Time, revokedAt time.Time) error
	RevokeFamily(ctx context.Context, tokenFamilyID string, revokedAt time.Time) error
	MarkReuseDetected(ctx context.Context, sessionID string, detectedAt time.Time) error
	UpdateLastSeen(ctx context.Context, sessionID string, lastSeenAt time.Time) error
}

type LoginSessionReader interface {
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*LoginSession, error)
	ListActiveByUserID(ctx context.Context, userID string, now time.Time) ([]LoginSession, error)
	CountActiveByUserID(ctx context.Context, userID string, now time.Time) (int, error)
}

type loginSessionReader struct {
	repository LoginSessionRepository
}

func NewLoginSessionReader(repository LoginSessionRepository) LoginSessionReader {
	return &loginSessionReader{repository: repository}
}

func (r *loginSessionReader) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*LoginSession, error) {
	return r.repository.FindByRefreshTokenHash(ctx, refreshTokenHash)
}

func (r *loginSessionReader) ListActiveByUserID(ctx context.Context, userID string, now time.Time) ([]LoginSession, error) {
	return r.repository.ListActiveByUserID(ctx, userID, now)
}

func (r *loginSessionReader) CountActiveByUserID(ctx context.Context, userID string, now time.Time) (int, error) {
	return r.repository.CountActiveByUserID(ctx, userID, now)
}

type LoginSessionWriter interface {
	Create(ctx context.Context, session LoginSession) error
	MarkRotated(ctx context.Context, sessionID string, rotatedAt time.Time) error
	Revoke(ctx context.Context, sessionID string, revokedAt time.Time) error
	RevokeByUserSessionID(ctx context.Context, userID string, sessionID string, revokedAt time.Time) error
	RevokeOldestActiveByUserID(ctx context.Context, userID string, now time.Time, revokedAt time.Time) error
	RevokeFamily(ctx context.Context, tokenFamilyID string, revokedAt time.Time) error
	MarkReuseDetected(ctx context.Context, sessionID string, detectedAt time.Time) error
	UpdateLastSeen(ctx context.Context, sessionID string, lastSeenAt time.Time) error
}

type loginSessionWriter struct {
	repository LoginSessionRepository
}

func NewLoginSessionWriter(repository LoginSessionRepository) LoginSessionWriter {
	return &loginSessionWriter{repository: repository}
}

func (w *loginSessionWriter) Create(ctx context.Context, session LoginSession) error {
	return w.repository.Create(ctx, session)
}

func (w *loginSessionWriter) MarkRotated(ctx context.Context, sessionID string, rotatedAt time.Time) error {
	return w.repository.MarkRotated(ctx, sessionID, rotatedAt)
}

func (w *loginSessionWriter) Revoke(ctx context.Context, sessionID string, revokedAt time.Time) error {
	return w.repository.Revoke(ctx, sessionID, revokedAt)
}

func (w *loginSessionWriter) RevokeByUserSessionID(ctx context.Context, userID string, sessionID string, revokedAt time.Time) error {
	return w.repository.RevokeByUserSessionID(ctx, userID, sessionID, revokedAt)
}

func (w *loginSessionWriter) RevokeOldestActiveByUserID(ctx context.Context, userID string, now time.Time, revokedAt time.Time) error {
	return w.repository.RevokeOldestActiveByUserID(ctx, userID, now, revokedAt)
}

func (w *loginSessionWriter) RevokeFamily(ctx context.Context, tokenFamilyID string, revokedAt time.Time) error {
	return w.repository.RevokeFamily(ctx, tokenFamilyID, revokedAt)
}

func (w *loginSessionWriter) MarkReuseDetected(ctx context.Context, sessionID string, detectedAt time.Time) error {
	return w.repository.MarkReuseDetected(ctx, sessionID, detectedAt)
}

func (w *loginSessionWriter) UpdateLastSeen(ctx context.Context, sessionID string, lastSeenAt time.Time) error {
	return w.repository.UpdateLastSeen(ctx, sessionID, lastSeenAt)
}
