package error

import (
	"errors"
	"log"
	stdhttp "net/http"

	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type APIError struct {
	Status  int
	Code    string
	Message string
	Level   coreerror.ErrorLevel
	Data    any
	Cause   error
}

func (e *APIError) Error() string {
	return e.Message
}

func New(apiErrorType APIErrorType, data any, cause error) *APIError {
	return &APIError{
		Status:  apiErrorType.Status,
		Code:    apiErrorType.Code,
		Message: apiErrorType.Message,
		Level:   apiErrorType.Level,
		Data:    data,
		Cause:   cause,
	}
}

func InvalidRequestError() *APIError {
	return New(InvalidRequest, nil, nil)
}

func NotFoundError() *APIError {
	return New(NotFound, nil, nil)
}

func Internal(cause error) *APIError {
	return New(DefaultError, nil, cause)
}

func FromError(err error) *APIError {
	if err == nil {
		return nil
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	var domainErr *coreerror.CoreException
	if errors.As(err, &domainErr) {
		return &APIError{
			Status:  mapKindToHTTPStatus(domainErr.Type.Kind),
			Code:    domainErr.Type.Code,
			Message: domainErr.Message,
			Level:   domainErr.Type.Level,
			Data:    domainErr.Data,
			Cause:   err,
		}
	}

	return Internal(err)
}

func Log(err *APIError) {
	if err == nil {
		return
	}

	switch err.Level {
	case coreerror.ErrorLevelError:
		log.Printf("ERROR code=%s message=%s cause=%v", err.Code, err.Message, err.Cause)
	case coreerror.ErrorLevelWarn:
		log.Printf("WARN code=%s message=%s cause=%v", err.Code, err.Message, err.Cause)
	default:
		log.Printf("INFO code=%s message=%s cause=%v", err.Code, err.Message, err.Cause)
	}
}

func mapKindToHTTPStatus(kind coreerror.ErrorKind) int {
	switch kind {
	case coreerror.ErrorKindUnauthorized:
		return stdhttp.StatusUnauthorized
	case coreerror.ErrorKindClient:
		return stdhttp.StatusBadRequest
	default:
		return stdhttp.StatusInternalServerError
	}
}
