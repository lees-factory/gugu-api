package auth

import "context"

type SessionRepository interface {
	Create(ctx context.Context, session Session) error
}
