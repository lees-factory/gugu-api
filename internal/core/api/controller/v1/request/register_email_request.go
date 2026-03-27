package request

import domainuser "github.com/ljj/gugu-api/internal/core/domain/user"

type RegisterEmail struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (r RegisterEmail) ToNewUser(passwordHash string) domainuser.NewUser {
	return domainuser.NewUser{
		Email:        r.Email,
		DisplayName:  r.DisplayName,
		PasswordHash: passwordHash,
		AuthSource:   "email",
	}
}
