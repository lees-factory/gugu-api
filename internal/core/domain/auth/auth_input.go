package auth

import "github.com/ljj/gugu-api/internal/core/domain/user"

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

type VerifyEmailInput struct {
	Token string
}

type LoginResult struct {
	User    user.User
	Session Session
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
