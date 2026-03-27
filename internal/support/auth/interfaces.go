package auth

import "time"

type IDGenerator interface {
	New() (string, error)
}

type TokenGenerator interface {
	New() (string, error)
}

type Clock interface {
	Now() time.Time
}

type AccessTokenIssuer interface {
	IssueAccessToken(userID string, now time.Time) (IssuedAccessToken, error)
}

type IssuedAccessToken struct {
	Token     string
	ExpiresAt time.Time
}

type PasswordVerifier interface {
	Verify(hashedPassword string, rawPassword string) error
}

type RefreshTokenHasher interface {
	Hash(value string) string
}
