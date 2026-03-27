package response

import "github.com/ljj/gugu-api/internal/core/domain/user"

type VerifyEmail struct {
	User User `json:"user"`
}

func NewVerifyEmail(source user.User) VerifyEmail {
	return VerifyEmail{
		User: NewUser(source),
	}
}
