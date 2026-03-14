package response

import "github.com/ljj/gugu-api/internal/core/domain/auth"

type Login struct {
	User   User   `json:"user"`
	Tokens Tokens `json:"tokens"`
}

func NewLogin(source auth.LoginResult) Login {
	return Login{
		User:   NewUser(source.User),
		Tokens: NewTokens(source.Tokens),
	}
}
