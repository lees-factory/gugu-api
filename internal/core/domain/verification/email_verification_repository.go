package verification

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, verification EmailVerification) error
	FindByCode(ctx context.Context, code string) (*EmailVerification, error)
	MarkUsed(ctx context.Context, code string, usedAt time.Time) error
}
