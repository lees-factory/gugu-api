package response

import "github.com/ljj/gugu-api/internal/core/domain/auth"

type RegisterEmail struct {
	User                   User   `json:"user"`
	VerificationToken      string `json:"verification_token,omitempty"`
	VerificationDispatched bool   `json:"verification_dispatched"`
}

func NewRegisterEmail(source auth.RegisterEmailResult) RegisterEmail {
	return RegisterEmail{
		User:                   NewUser(source.User),
		VerificationToken:      source.VerificationToken,
		VerificationDispatched: source.VerificationDispatched,
	}
}
