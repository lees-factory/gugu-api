package auth

import (
	"context"
	"sync"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

type SessionMemoryRepository struct {
	mu       sync.RWMutex
	sessions map[string]domainauth.Session
}

func NewSessionRepository() *SessionMemoryRepository {
	return &SessionMemoryRepository{
		sessions: make(map[string]domainauth.Session),
	}
}

func (r *SessionMemoryRepository) Create(_ context.Context, session domainauth.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.Token] = session
	return nil
}
