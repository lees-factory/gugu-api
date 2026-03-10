package verification

import "time"

type EmailVerification struct {
	Token     string
	UserID    string
	Email     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
