package response

import (
	"time"

	supportauth "github.com/ljj/gugu-api/internal/support/auth"
)

type Tokens struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func NewTokens(source supportauth.AuthTokens) Tokens {
	return Tokens{
		AccessToken:      source.AccessToken,
		RefreshToken:     source.RefreshToken,
		TokenType:        source.TokenType,
		AccessExpiresAt:  source.AccessExpiresAt,
		RefreshExpiresAt: source.RefreshExpiresAt,
	}
}

func NewTokensFromSource(source *supportauth.AuthTokens) Tokens {
	return NewTokens(*source)
}
