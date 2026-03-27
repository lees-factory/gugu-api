package response

import (
	apierror "github.com/ljj/gugu-api/internal/core/api/support/error"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewAPIErrorMessage(errorType apierror.APIErrorType, data any) *ErrorMessage {
	return &ErrorMessage{
		Code:    errorType.Code,
		Message: errorType.Message,
		Data:    data,
	}
}

func NewCoreErrorMessage(errorType coreerror.ErrorType, data any) *ErrorMessage {
	return &ErrorMessage{
		Code:    errorType.Code,
		Message: errorType.Message,
		Data:    data,
	}
}
