package coreerror

type ErrorType struct {
	Kind    ErrorKind
	Code    string
	Message string
	Level   ErrorLevel
}

// User (G1xxx)
var (
	EmailAlreadyExists   = ErrorType{KindClient, G1000, "email already exists", LevelInfo}
	VerificationNotFound = ErrorType{KindClient, G2002, "verification code not found", LevelInfo}
)

// Auth (G2xxx)
var (
	InvalidCredentials   = ErrorType{KindUnauthorized, G2000, "invalid credentials", LevelWarn}
	EmailNotVerified     = ErrorType{KindClient, G2001, "email is not verified", LevelInfo}
	OAuthProviderInvalid = ErrorType{KindClient, G2003, "oauth provider is invalid", LevelInfo}
	RefreshTokenInvalid  = ErrorType{KindUnauthorized, G2004, "refresh token is invalid", LevelWarn}
)

// Product & TrackedItem (G3xxx)
var (
	UnsupportedMarket   = ErrorType{KindClient, G3000, "unsupported market", LevelInfo}
	TrackedItemExists   = ErrorType{KindClient, G3001, "tracked item already exists", LevelInfo}
	ProductNotFound     = ErrorType{KindClient, G3002, "product not found", LevelInfo}
	TrackedItemNotFound = ErrorType{KindClient, G3003, "tracked item not found", LevelInfo}
)
