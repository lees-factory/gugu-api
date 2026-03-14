package auth

import "time"

type AuthTokenIssuer interface {
	Issue(userID string, now time.Time) (AuthTokens, error)
}
