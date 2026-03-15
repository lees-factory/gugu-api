package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type OAuthIdentitySQLCRepository struct {
	queries *sqldb.Queries
}

func NewOAuthIdentitySQLCRepository(db *sql.DB) *OAuthIdentitySQLCRepository {
	return &OAuthIdentitySQLCRepository{queries: sqldb.New(db)}
}

func (r *OAuthIdentitySQLCRepository) FindByProviderSubject(ctx context.Context, provider string, subject string) (*domainauth.OAuthIdentity, error) {
	row, err := r.queries.FindOAuthIdentity(ctx, sqldb.FindOAuthIdentityParams{
		Provider: provider,
		Subject:  subject,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domainauth.OAuthIdentity{
		ID:          row.ID,
		UserID:      row.UserID,
		Provider:    row.Provider,
		Subject:     row.Subject,
		Email:       row.Email,
		CreatedAt:   row.CreatedAt,
		LastLoginAt: row.LastLoginAt,
	}, nil
}

func (r *OAuthIdentitySQLCRepository) Create(ctx context.Context, identity domainauth.OAuthIdentity) error {
	return r.queries.CreateOAuthIdentity(ctx, sqldb.CreateOAuthIdentityParams{
		ID:          identity.ID,
		UserID:      identity.UserID,
		Provider:    identity.Provider,
		Subject:     identity.Subject,
		Email:       identity.Email,
		CreatedAt:   identity.CreatedAt,
		LastLoginAt: identity.LastLoginAt,
	})
}

func (r *OAuthIdentitySQLCRepository) UpdateLastLogin(ctx context.Context, provider string, subject string, lastLoginAt time.Time) error {
	affected, err := r.queries.UpdateOAuthLastLogin(ctx, sqldb.UpdateOAuthLastLoginParams{
		Provider:    provider,
		Subject:     subject,
		LastLoginAt: lastLoginAt,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
