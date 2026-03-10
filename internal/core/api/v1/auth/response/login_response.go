package response

import "github.com/ljj/gugu-api/internal/core/domain/auth"

type Login struct {
	User    User    `json:"user"`
	Session Session `json:"session"`
}

func NewLogin(source auth.LoginResult) Login {
	return Login{
		User:    NewUser(source.User),
		Session: NewSession(source.Session),
	}
}
