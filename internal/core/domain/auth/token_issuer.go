package auth

import "time"

type AccessTokenIssuer interface {
	IssueAccessToken(userID string, now time.Time) (IssuedAccessToken, error)
}

type IssuedAccessToken struct {
	Token     string
	ExpiresAt time.Time
}
