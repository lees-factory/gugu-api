package auth

import (
	"context"
	"time"
)

type OAuthIdentityRepository interface {
	FindByProviderSubject(ctx context.Context, provider string, subject string) (*OAuthIdentity, error)
	Create(ctx context.Context, identity OAuthIdentity) error
	UpdateLastLogin(ctx context.Context, provider string, subject string, lastLoginAt time.Time) error
}
