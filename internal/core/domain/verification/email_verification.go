package verification

import "time"

type EmailVerification struct {
	Code      string
	UserID    string
	Email     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
