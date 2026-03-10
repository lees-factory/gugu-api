package user

import (
	"context"
	"time"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userID string) (*User, error)
	Create(ctx context.Context, user User) error
	MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) error
}
