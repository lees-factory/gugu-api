package auth

import domainuser "github.com/ljj/gugu-api/internal/core/domain/user"

type LoginEmailInput struct {
	Email      string
	Password   string
	UserAgent  string
	ClientIP   string
	DeviceName string
}

type LoginResult struct {
	User   domainuser.User
	Tokens AuthTokens
}

type OAuthLoginInput struct {
	Provider    OAuthProvider
	Subject     string
	Email       string
	DisplayName string
	UserAgent   string
	ClientIP    string
	DeviceName  string
}

type SessionMetadata struct {
	UserAgent  string
	ClientIP   string
	DeviceName string
}

type RefreshTokensInput struct {
	RefreshToken string
	UserAgent    string
	ClientIP     string
	DeviceName   string
}

type LogoutInput struct {
	RefreshToken string
}

type CreateOAuthIdentityInput struct {
	UserID   string
	Provider OAuthProvider
	Subject  string
	Email    string
}
