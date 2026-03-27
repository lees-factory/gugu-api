package auth

import (
	"context"
	"sync"
	"time"

	supportauth "github.com/ljj/gugu-api/internal/support/auth"
)

type LoginSessionMemoryRepository struct {
	mu                 sync.RWMutex
	sessions           map[string]supportauth.LoginSession
	sessionIDsByHash   map[string]string
	sessionIDsByFamily map[string]map[string]struct{}
}

func NewLoginSessionRepository() *LoginSessionMemoryRepository {
	return &LoginSessionMemoryRepository{
		sessions:           make(map[string]supportauth.LoginSession),
		sessionIDsByHash:   make(map[string]string),
		sessionIDsByFamily: make(map[string]map[string]struct{}),
	}
}

func (r *LoginSessionMemoryRepository) Create(_ context.Context, session supportauth.LoginSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	r.sessionIDsByHash[session.RefreshTokenHash] = session.ID
	if _, ok := r.sessionIDsByFamily[session.TokenFamilyID]; !ok {
		r.sessionIDsByFamily[session.TokenFamilyID] = make(map[string]struct{})
	}
	r.sessionIDsByFamily[session.TokenFamilyID][session.ID] = struct{}{}
	return nil
}

func (r *LoginSessionMemoryRepository) FindByRefreshTokenHash(_ context.Context, refreshTokenHash string) (*supportauth.LoginSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sessionID, ok := r.sessionIDsByHash[refreshTokenHash]
	if !ok {
		return nil, nil
	}

	session := r.sessions[sessionID]
	return &session, nil
}

func (r *LoginSessionMemoryRepository) MarkRotated(_ context.Context, sessionID string, rotatedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session := r.sessions[sessionID]
	session.RotatedAt = &rotatedAt
	r.sessions[sessionID] = session
	return nil
}

func (r *LoginSessionMemoryRepository) Revoke(_ context.Context, sessionID string, revokedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session := r.sessions[sessionID]
	session.RevokedAt = &revokedAt
	r.sessions[sessionID] = session
	return nil
}

func (r *LoginSessionMemoryRepository) RevokeFamily(_ context.Context, tokenFamilyID string, revokedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for sessionID := range r.sessionIDsByFamily[tokenFamilyID] {
		session := r.sessions[sessionID]
		session.RevokedAt = &revokedAt
		r.sessions[sessionID] = session
	}
	return nil
}

func (r *LoginSessionMemoryRepository) MarkReuseDetected(_ context.Context, sessionID string, detectedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session := r.sessions[sessionID]
	session.ReuseDetectedAt = &detectedAt
	r.sessions[sessionID] = session
	return nil
}

func (r *LoginSessionMemoryRepository) UpdateLastSeen(_ context.Context, sessionID string, lastSeenAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session := r.sessions[sessionID]
	session.LastSeenAt = lastSeenAt
	r.sessions[sessionID] = session
	return nil
}
