package auth

import (
	"context"
	"time"
)

type LoginSessionRepository interface {
	Create(ctx context.Context, session LoginSession) error
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*LoginSession, error)
	MarkRotated(ctx context.Context, sessionID string, rotatedAt time.Time) error
	Revoke(ctx context.Context, sessionID string, revokedAt time.Time) error
	RevokeFamily(ctx context.Context, tokenFamilyID string, revokedAt time.Time) error
	MarkReuseDetected(ctx context.Context, sessionID string, detectedAt time.Time) error
	UpdateLastSeen(ctx context.Context, sessionID string, lastSeenAt time.Time) error
}
