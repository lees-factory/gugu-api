package response

import "github.com/ljj/gugu-api/internal/core/domain/user"

type RegisterEmail struct {
	User                   User   `json:"user"`
	VerificationToken      string `json:"verification_token,omitempty"`
	VerificationDispatched bool   `json:"verification_dispatched"`
}

func NewRegisterEmail(user user.User, verificationToken string, verificationDispatched bool) RegisterEmail {
	return RegisterEmail{
		User:                   NewUser(user),
		VerificationToken:      verificationToken,
		VerificationDispatched: verificationDispatched,
	}
}
