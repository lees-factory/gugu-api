package error

import (
	stdhttp "net/http"

	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type APIErrorType struct {
	Status  int
	Code    string
	Message string
	Level   coreerror.ErrorLevel
}

func (e APIErrorType) CodeValue() string {
	return e.Code
}

func (e APIErrorType) MessageValue() string {
	return e.Message
}

var (
	DefaultError = APIErrorType{
		Status:  stdhttp.StatusInternalServerError,
		Code:    "E500",
		Message: "An unexpected error has occurred.",
		Level:   coreerror.LevelError,
	}
	NotFound = APIErrorType{
		Status:  stdhttp.StatusNotFound,
		Code:    "E501",
		Message: "Not Found",
		Level:   coreerror.LevelInfo,
	}
	InvalidRequest = APIErrorType{
		Status:  stdhttp.StatusBadRequest,
		Code:    "E400",
		Message: "invalid request body",
		Level:   coreerror.LevelInfo,
	}
)
