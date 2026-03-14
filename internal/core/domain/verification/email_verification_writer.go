package verification

import (
	"context"
	"time"
)

type Writer interface {
	Create(ctx context.Context, emailVerification EmailVerification) error
	MarkUsed(ctx context.Context, code string, usedAt time.Time) error
}

type writer struct {
	repository Repository
}

func NewWriter(repository Repository) Writer {
	return &writer{repository: repository}
}

func (w *writer) Create(ctx context.Context, emailVerification EmailVerification) error {
	return w.repository.Create(ctx, emailVerification)
}

func (w *writer) MarkUsed(ctx context.Context, code string, usedAt time.Time) error {
	return w.repository.MarkUsed(ctx, code, usedAt)
}
