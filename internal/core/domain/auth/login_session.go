package auth

import "time"

type LoginSession struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	TokenFamilyID    string
	ParentSessionID  *string
	UserAgent        string
	ClientIP         string
	DeviceName       string
	ExpiresAt        time.Time
	LastSeenAt       time.Time
	RotatedAt        *time.Time
	RevokedAt        *time.Time
	ReuseDetectedAt  *time.Time
	CreatedAt        time.Time
}
