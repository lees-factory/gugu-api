package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

type OAuthIdentitySQLCRepository struct {
	db *sql.DB
}

func NewOAuthIdentityRepository(db *sql.DB) *OAuthIdentitySQLCRepository {
	return &OAuthIdentitySQLCRepository{db: db}
}

func (r *OAuthIdentitySQLCRepository) FindByProviderSubject(ctx context.Context, provider string, subject string) (*domainauth.OAuthIdentity, error) {
	const query = `
SELECT
	id,
	user_id,
	provider,
	subject,
	email,
	created_at,
	last_login_at
	FROM gugu.oauth_identities
WHERE provider = $1 AND subject = $2
`

	var identity domainauth.OAuthIdentity
	err := r.db.QueryRowContext(ctx, query, provider, subject).Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.Subject,
		&identity.Email,
		&identity.CreatedAt,
		&identity.LastLoginAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &identity, nil
}

func (r *OAuthIdentitySQLCRepository) Create(ctx context.Context, identity domainauth.OAuthIdentity) error {
	const query = `
	INSERT INTO gugu.oauth_identities (
	id,
	user_id,
	provider,
	subject,
	email,
	created_at,
	last_login_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
)
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		identity.ID,
		identity.UserID,
		identity.Provider,
		identity.Subject,
		identity.Email,
		identity.CreatedAt,
		identity.LastLoginAt,
	)
	return err
}

func (r *OAuthIdentitySQLCRepository) UpdateLastLogin(ctx context.Context, provider string, subject string, lastLoginAt time.Time) error {
	const query = `
	UPDATE gugu.oauth_identities
SET last_login_at = $3
WHERE provider = $1 AND subject = $2
`

	result, err := r.db.ExecContext(ctx, query, provider, subject, lastLoginAt)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
