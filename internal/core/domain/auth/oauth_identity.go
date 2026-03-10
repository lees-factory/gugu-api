package auth

import "time"

type OAuthIdentity struct {
	ID          string
	UserID      string
	Provider    string
	Subject     string
	Email       string
	CreatedAt   time.Time
	LastLoginAt time.Time
}
