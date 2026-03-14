package verification

import (
	"context"
	"errors"
	"time"

	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
)

var ErrNotImplemented = errors.New("email verification sqlc repository is not implemented yet")

type EmailVerificationSQLCRepository struct{}

func NewRepository() *EmailVerificationSQLCRepository {
	return &EmailVerificationSQLCRepository{}
}

func (r *EmailVerificationSQLCRepository) Create(context.Context, domainverification.EmailVerification) error {
	return ErrNotImplemented
}

func (r *EmailVerificationSQLCRepository) FindByCode(context.Context, string) (*domainverification.EmailVerification, error) {
	return nil, ErrNotImplemented
}

func (r *EmailVerificationSQLCRepository) MarkUsed(context.Context, string, time.Time) error {
	return ErrNotImplemented
}
