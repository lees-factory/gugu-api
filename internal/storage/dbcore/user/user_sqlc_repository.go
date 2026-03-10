package user

import (
	"context"
	"errors"
	"time"

	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
)

var ErrNotImplemented = errors.New("user sqlc repository is not implemented yet")

type UserSQLCRepository struct{}

func NewRepository() *UserSQLCRepository {
	return &UserSQLCRepository{}
}

func (r *UserSQLCRepository) FindByEmail(context.Context, string) (*domainuser.User, error) {
	return nil, ErrNotImplemented
}

func (r *UserSQLCRepository) FindByID(context.Context, string) (*domainuser.User, error) {
	return nil, ErrNotImplemented
}

func (r *UserSQLCRepository) Create(context.Context, domainuser.User) error {
	return ErrNotImplemented
}

func (r *UserSQLCRepository) MarkEmailVerified(context.Context, string, time.Time) error {
	return ErrNotImplemented
}
