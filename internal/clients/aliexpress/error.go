package aliexpress

import "fmt"

type RemoteError struct {
	Code      string
	Message   string
	RequestID string
}

func (e *RemoteError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprintf("aliexpress error: code=%s message=%s", e.Code, e.Message)
}
