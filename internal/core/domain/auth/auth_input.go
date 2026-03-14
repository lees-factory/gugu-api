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
	VerificationCode       string
	VerificationDispatched bool
}

type RegisterEmailInput struct {
	Email       string
	Password    string
	DisplayName string
}

type RegisterEmailResult struct {
	User                   user.User
	VerificationCode       string
	VerificationDispatched bool
}

type LoginEmailInput struct {
	Email      string
	Password   string
	UserAgent  string
	ClientIP   string
	DeviceName string
}

type LoginResult struct {
	User   user.User
	Tokens AuthTokens
}

type VerifyEmailCodeInput struct {
	Code string
}

type VerifyEmailCodeResult struct {
	UserID     string
	VerifiedAt time.Time
}

type VerifyEmailInput struct {
	Code string
}

type VerifyEmailResult struct {
	User user.User
}

type OAuthLoginInput struct {
	Provider    Provider
	Subject     string
	Email       string
	DisplayName string
	UserAgent   string
	ClientIP    string
	DeviceName  string
}
