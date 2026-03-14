package auth

import coreerror "github.com/ljj/gugu-api/internal/core/error"

var (
	ErrorTypeInvalidCredentials = coreerror.ErrorType{
		Kind:    coreerror.ErrorKindUnauthorized,
		Code:    "B2000",
		Message: "invalid credentials",
		Level:   coreerror.ErrorLevelWarn,
	}
	ErrorTypeEmailNotVerified = coreerror.ErrorType{
		Kind:    coreerror.ErrorKindClient,
		Code:    "B2001",
		Message: "email is not verified",
		Level:   coreerror.ErrorLevelInfo,
	}
	ErrorTypeVerificationNotFound = coreerror.ErrorType{
		Kind:    coreerror.ErrorKindClient,
		Code:    "B2002",
		Message: "verification code not found",
		Level:   coreerror.ErrorLevelInfo,
	}
	ErrorTypeOAuthProviderInvalid = coreerror.ErrorType{
		Kind:    coreerror.ErrorKindClient,
		Code:    "B2003",
		Message: "oauth provider is invalid",
		Level:   coreerror.ErrorLevelInfo,
	}
)

var (
	ErrInvalidCredentials   = coreerror.New(ErrorTypeInvalidCredentials)
	ErrEmailNotVerified     = coreerror.New(ErrorTypeEmailNotVerified)
	ErrVerificationNotFound = coreerror.New(ErrorTypeVerificationNotFound)
	ErrOAuthProviderInvalid = coreerror.New(ErrorTypeOAuthProviderInvalid)
)
