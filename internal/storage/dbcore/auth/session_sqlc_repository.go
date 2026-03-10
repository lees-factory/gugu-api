package auth

import (
	"context"
	"errors"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

var ErrSessionNotImplemented = errors.New("session sqlc repository is not implemented yet")

type SessionSQLCRepository struct{}

func NewSessionRepository() *SessionSQLCRepository {
	return &SessionSQLCRepository{}
}

func (r *SessionSQLCRepository) Create(context.Context, domainauth.Session) error {
	return ErrSessionNotImplemented
}
