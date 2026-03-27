package response

import supportauth "github.com/ljj/gugu-api/internal/support/auth"

type Login struct {
	User   User   `json:"user"`
	Tokens Tokens `json:"tokens"`
}

func NewLogin(source supportauth.LoginResult) Login {
	return Login{
		User:   NewUser(source.User),
		Tokens: NewTokens(source.Tokens),
	}
}
