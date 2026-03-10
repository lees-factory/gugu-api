package verification

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, verification EmailVerification) error
	FindByToken(ctx context.Context, token string) (*EmailVerification, error)
	MarkUsed(ctx context.Context, token string, usedAt time.Time) error
}
