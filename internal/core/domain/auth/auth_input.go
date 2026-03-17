package auth

import "github.com/ljj/gugu-api/internal/core/domain/user"

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

type OAuthLoginInput struct {
	Provider    Provider
	Subject     string
	Email       string
	DisplayName string
	UserAgent   string
	ClientIP    string
	DeviceName  string
}
