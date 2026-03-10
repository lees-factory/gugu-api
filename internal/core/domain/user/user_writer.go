package user

import (
	"context"
	"time"
)

type Writer interface {
	Create(ctx context.Context, user User) error
	MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) error
}

type writer struct {
	repository Repository
}

func NewWriter(repository Repository) Writer {
	return &writer{repository: repository}
}

func (w *writer) Create(ctx context.Context, user User) error {
	return w.repository.Create(ctx, user)
}

func (w *writer) MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) error {
	return w.repository.MarkEmailVerified(ctx, userID, verifiedAt)
}
