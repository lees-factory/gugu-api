package auth

import (
	"context"
	"time"
)

type IDGenerator interface {
	New() (string, error)
}

type TokenGenerator interface {
	New() (string, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword string, rawPassword string) error
}

type Clock interface {
	Now() time.Time
}

type VerificationSender interface {
	SendVerification(ctx context.Context, email string, token string) error
}
