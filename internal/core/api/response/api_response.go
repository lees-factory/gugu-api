package response

type ResultType string

const (
	ResultTypeSuccess ResultType = "SUCCESS"
	ResultTypeError   ResultType = "ERROR"
)

type APIResponse[T any] struct {
	Result ResultType    `json:"result"`
	Data   *T            `json:"data,omitempty"`
	Error  *ErrorMessage `json:"error,omitempty"`
}

func Success() APIResponse[any] {
	return APIResponse[any]{
		Result: ResultTypeSuccess,
	}
}

func SuccessWithData[T any](data T) APIResponse[T] {
	return APIResponse[T]{
		Result: ResultTypeSuccess,
		Data:   &data,
	}
}

func ErrorFromAPIError[T any](errorType interface {
	CodeValue() string
	MessageValue() string
}, data any) APIResponse[T] {
	return APIResponse[T]{
		Result: ResultTypeError,
		Error: &ErrorMessage{
			Code:    errorType.CodeValue(),
			Message: errorType.MessageValue(),
			Data:    data,
		},
	}
}

func ErrorFromMessage[T any](error *ErrorMessage) APIResponse[T] {
	return APIResponse[T]{
		Result: ResultTypeError,
		Error:  error,
	}
}
