package coreerror

type ErrorType struct {
	Kind    ErrorKind
	Code    string
	Message string
	Level   ErrorLevel
}

func (e ErrorType) CodeValue() string {
	return e.Code
}

func (e ErrorType) MessageValue() string {
	return e.Message
}
