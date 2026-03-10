package response

import "github.com/ljj/gugu-api/internal/core/domain/auth"

type VerifyEmail struct {
	User User `json:"user"`
}

func NewVerifyEmail(source auth.VerifyEmailResult) VerifyEmail {
	return VerifyEmail{
		User: NewUser(source.User),
	}
}
