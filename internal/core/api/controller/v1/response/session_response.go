package response

import (
	"time"

	supportauth "github.com/ljj/gugu-api/internal/support/auth"
)

type LoginSession struct {
	ID         string    `json:"id"`
	UserAgent  string    `json:"user_agent"`
	ClientIP   string    `json:"client_ip"`
	DeviceName string    `json:"device_name"`
	ExpiresAt  time.Time `json:"expires_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewLoginSessions(source []supportauth.LoginSession) []LoginSession {
	items := make([]LoginSession, 0, len(source))
	for _, session := range source {
		items = append(items, LoginSession{
			ID:         session.ID,
			UserAgent:  session.UserAgent,
			ClientIP:   session.ClientIP,
			DeviceName: session.DeviceName,
			ExpiresAt:  session.ExpiresAt,
			LastSeenAt: session.LastSeenAt,
			CreatedAt:  session.CreatedAt,
		})
	}
	return items
}
