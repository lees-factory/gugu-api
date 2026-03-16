package product

import coreerror "github.com/ljj/gugu-api/internal/core/error"

var ErrorTypeUnsupportedMarket = coreerror.ErrorType{
	Kind:    coreerror.ErrorKindClient,
	Code:    "B3000",
	Message: "unsupported market",
	Level:   coreerror.ErrorLevelInfo,
}

var ErrorTypeProductNotFound = coreerror.ErrorType{
	Kind:    coreerror.ErrorKindClient,
	Code:    "B3002",
	Message: "product not found",
	Level:   coreerror.ErrorLevelInfo,
}

var ErrUnsupportedMarket = coreerror.New(ErrorTypeUnsupportedMarket)
var ErrProductNotFound = coreerror.New(ErrorTypeProductNotFound)
