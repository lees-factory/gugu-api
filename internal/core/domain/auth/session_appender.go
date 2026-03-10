package auth

import "context"

type SessionAppender interface {
	Create(ctx context.Context, session Session) error
}

type sessionAppender struct {
	repository SessionRepository
}

func NewSessionAppender(repository SessionRepository) SessionAppender {
	return &sessionAppender{repository: repository}
}

func (w *sessionAppender) Create(ctx context.Context, session Session) error {
	return w.repository.Create(ctx, session)
}
