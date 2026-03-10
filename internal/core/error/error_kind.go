package coreerror

type ErrorKind string

const (
	ErrorKindClient       ErrorKind = "CLIENT_ERROR"
	ErrorKindServer       ErrorKind = "SERVER_ERROR"
	ErrorKindUnauthorized ErrorKind = "UNAUTHORIZED"
)
