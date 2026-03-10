package response

import (
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/auth"
)

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewSession(source auth.Session) Session {
	return Session{
		Token:     source.Token,
		UserID:    source.UserID,
		ExpiresAt: source.ExpiresAt,
	}
}
