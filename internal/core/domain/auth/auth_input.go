package auth

import (
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/user"
)

type IssueEmailVerificationInput struct {
	UserID string
	Email  string
}

type IssueEmailVerificationResult struct {
	VerificationToken      string
	VerificationDispatched bool
}

type RegisterEmailInput struct {
	Email       string
	Password    string
	DisplayName string
}

type RegisterEmailResult struct {
	User                   user.User
	VerificationToken      string
	VerificationDispatched bool
}

type LoginEmailInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	User    user.User
	Session Session
}

type VerifyEmailTokenInput struct {
	Token string
}

type VerifyEmailTokenResult struct {
	UserID     string
	VerifiedAt time.Time
}

type VerifyEmailInput struct {
	Token string
}

type VerifyEmailResult struct {
	User user.User
}

type OAuthLoginInput struct {
	Provider    Provider
	Subject     string
	Email       string
	DisplayName string
}
