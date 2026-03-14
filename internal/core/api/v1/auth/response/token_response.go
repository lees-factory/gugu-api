package response

import (
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/auth"
)

type Tokens struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func NewTokens(source auth.AuthTokens) Tokens {
	return Tokens{
		AccessToken:      source.AccessToken,
		RefreshToken:     source.RefreshToken,
		TokenType:        source.TokenType,
		AccessExpiresAt:  source.AccessExpiresAt,
		RefreshExpiresAt: source.RefreshExpiresAt,
	}
}

func NewTokensFromSource(source *auth.AuthTokens) Tokens {
	return NewTokens(*source)
}
