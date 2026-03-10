package auth

import (
	"context"
	"errors"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

var ErrOAuthNotImplemented = errors.New("oauth identity sqlc repository is not implemented yet")

type OAuthIdentitySQLCRepository struct{}

func NewOAuthIdentityRepository() *OAuthIdentitySQLCRepository {
	return &OAuthIdentitySQLCRepository{}
}

func (r *OAuthIdentitySQLCRepository) FindByProviderSubject(context.Context, string, string) (*domainauth.OAuthIdentity, error) {
	return nil, ErrOAuthNotImplemented
}

func (r *OAuthIdentitySQLCRepository) Create(context.Context, domainauth.OAuthIdentity) error {
	return ErrOAuthNotImplemented
}

func (r *OAuthIdentitySQLCRepository) UpdateLastLogin(context.Context, string, string, time.Time) error {
	return ErrOAuthNotImplemented
}
