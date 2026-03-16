package trackeditem

import coreerror "github.com/ljj/gugu-api/internal/core/error"

var ErrorTypeAlreadyExists = coreerror.ErrorType{
	Kind:    coreerror.ErrorKindClient,
	Code:    "B3001",
	Message: "tracked item already exists",
	Level:   coreerror.ErrorLevelInfo,
}

var ErrorTypeTrackedItemNotFound = coreerror.ErrorType{
	Kind:    coreerror.ErrorKindClient,
	Code:    "B3003",
	Message: "tracked item not found",
	Level:   coreerror.ErrorLevelInfo,
}

var ErrAlreadyExists = coreerror.New(ErrorTypeAlreadyExists)
var ErrTrackedItemNotFound = coreerror.New(ErrorTypeTrackedItemNotFound)
