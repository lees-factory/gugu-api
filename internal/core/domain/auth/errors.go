package auth

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrEmailNotVerified     = errors.New("email is not verified")
	ErrVerificationNotFound = errors.New("verification token not found")
	ErrOAuthProviderInvalid = errors.New("oauth provider is invalid")
)
