package verification

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
)

type EmailVerificationSQLCRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *EmailVerificationSQLCRepository {
	return &EmailVerificationSQLCRepository{db: db}
}

func (r *EmailVerificationSQLCRepository) Create(ctx context.Context, emailVerification domainverification.EmailVerification) error {
	const query = `
	INSERT INTO gugu.email_verifications (
	code,
	user_id,
	email,
	expires_at,
	used_at,
	created_at
) VALUES (
	$1, $2, $3, $4, $5, $6
)
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		emailVerification.Code,
		emailVerification.UserID,
		emailVerification.Email,
		emailVerification.ExpiresAt,
		emailVerification.UsedAt,
		emailVerification.CreatedAt,
	)
	return err
}

func (r *EmailVerificationSQLCRepository) FindByCode(ctx context.Context, code string) (*domainverification.EmailVerification, error) {
	const query = `
SELECT
	code,
	user_id,
	email,
	expires_at,
	used_at,
	created_at
	FROM gugu.email_verifications
WHERE code = $1
`

	var verification domainverification.EmailVerification
	var usedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&verification.Code,
		&verification.UserID,
		&verification.Email,
		&verification.ExpiresAt,
		&usedAt,
		&verification.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if usedAt.Valid {
		verification.UsedAt = &usedAt.Time
	}

	return &verification, nil
}

func (r *EmailVerificationSQLCRepository) MarkUsed(ctx context.Context, code string, usedAt time.Time) error {
	const query = `
	UPDATE gugu.email_verifications
SET used_at = $2
WHERE code = $1
`

	result, err := r.db.ExecContext(ctx, query, code, usedAt)
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
