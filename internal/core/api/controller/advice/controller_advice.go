package advice

import (
	"encoding/json"
	"fmt"
	stdhttp "net/http"

	apierror "github.com/ljj/gugu-api/internal/core/support/error"
	"github.com/ljj/gugu-api/internal/core/support/response"
)

type HandlerFunc func(r *stdhttp.Request) (int, any, error)

func Wrap(handler HandlerFunc) stdhttp.HandlerFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		status, payload, err := handler(r)
		if err != nil {
			writeAPIError(w, apierror.FromError(err))
			return
		}
		_ = writeJSON(w, status, payload)
	}
}

func DecodeJSON(r *stdhttp.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return apierror.InvalidRequestError()
	}
	return nil
}

func writeAPIError(w stdhttp.ResponseWriter, err *apierror.APIError) {
	if err == nil {
		err = apierror.Internal(nil)
	}
	apierror.Log(err)
	_ = writeJSON(w, err.Status, response.ErrorFromMessage[any](&response.ErrorMessage{
		Code:    err.Code,
		Message: err.Message,
		Data:    err.Data,
	}))
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}
