package auth

import "context"

type LoginSessionReader interface {
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*LoginSession, error)
}

type loginSessionReader struct {
	repository LoginSessionRepository
}

func NewLoginSessionReader(repository LoginSessionRepository) LoginSessionReader {
	return &loginSessionReader{repository: repository}
}

func (r *loginSessionReader) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*LoginSession, error) {
	return r.repository.FindByRefreshTokenHash(ctx, refreshTokenHash)
}
