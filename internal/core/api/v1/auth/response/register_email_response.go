package response

import "github.com/ljj/gugu-api/internal/core/domain/user"

type RegisterEmail struct {
	User                   User   `json:"user"`
	VerificationCode       string `json:"verification_code,omitempty"`
	VerificationDispatched bool   `json:"verification_dispatched"`
}

func NewRegisterEmail(user user.User, verificationCode string, verificationDispatched bool) RegisterEmail {
	return RegisterEmail{
		User:                   NewUser(user),
		VerificationCode:       verificationCode,
		VerificationDispatched: verificationDispatched,
	}
}
