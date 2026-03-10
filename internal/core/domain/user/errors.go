package user

import coreerror "github.com/ljj/gugu-api/internal/core/error"

var ErrorTypeEmailAlreadyExists = coreerror.ErrorType{
	Kind:    coreerror.ErrorKindClient,
	Code:    "B1000",
	Message: "email already exists",
	Level:   coreerror.ErrorLevelInfo,
}

var ErrEmailAlreadyExists = coreerror.New(ErrorTypeEmailAlreadyExists)
