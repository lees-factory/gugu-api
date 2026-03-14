package auth

import "time"

type AuthTokens struct {
	AccessToken      string
	RefreshToken     string
	TokenType        string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}
