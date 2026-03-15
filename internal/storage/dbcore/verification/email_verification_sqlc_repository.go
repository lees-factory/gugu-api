package verification

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type EmailVerificationSQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *EmailVerificationSQLCRepository {
	return &EmailVerificationSQLCRepository{queries: sqldb.New(db)}
}

func (r *EmailVerificationSQLCRepository) Create(ctx context.Context, emailVerification domainverification.EmailVerification) error {
	return r.queries.CreateEmailVerification(ctx, sqldb.CreateEmailVerificationParams{
		Code:      emailVerification.Code,
		UserID:    emailVerification.UserID,
		Email:     emailVerification.Email,
		ExpiresAt: emailVerification.ExpiresAt,
		UsedAt:    nullTime(emailVerification.UsedAt),
		CreatedAt: emailVerification.CreatedAt,
	})
}

func (r *EmailVerificationSQLCRepository) FindByCode(ctx context.Context, code string) (*domainverification.EmailVerification, error) {
	row, err := r.queries.FindEmailVerificationByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	verification := &domainverification.EmailVerification{
		Code:      row.Code,
		UserID:    row.UserID,
		Email:     row.Email,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}
	if row.UsedAt.Valid {
		verification.UsedAt = &row.UsedAt.Time
	}

	return verification, nil
}

func (r *EmailVerificationSQLCRepository) MarkUsed(ctx context.Context, code string, usedAt time.Time) error {
	affected, err := r.queries.MarkEmailVerificationUsed(ctx, sqldb.MarkEmailVerificationUsedParams{
		Code:   code,
		UsedAt: sql.NullTime{Time: usedAt, Valid: true},
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func nullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *value, Valid: true}
}
