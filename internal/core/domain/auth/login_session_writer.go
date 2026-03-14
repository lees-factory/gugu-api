package auth

import (
	"context"
	"time"
)

type LoginSessionWriter interface {
	Create(ctx context.Context, session LoginSession) error
	MarkRotated(ctx context.Context, sessionID string, rotatedAt time.Time) error
	Revoke(ctx context.Context, sessionID string, revokedAt time.Time) error
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

func (w *loginSessionWriter) RevokeFamily(ctx context.Context, tokenFamilyID string, revokedAt time.Time) error {
	return w.repository.RevokeFamily(ctx, tokenFamilyID, revokedAt)
}

func (w *loginSessionWriter) MarkReuseDetected(ctx context.Context, sessionID string, detectedAt time.Time) error {
	return w.repository.MarkReuseDetected(ctx, sessionID, detectedAt)
}

func (w *loginSessionWriter) UpdateLastSeen(ctx context.Context, sessionID string, lastSeenAt time.Time) error {
	return w.repository.UpdateLastSeen(ctx, sessionID, lastSeenAt)
}
