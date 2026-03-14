package auth

type SessionMetadata struct {
	UserAgent  string
	ClientIP   string
	DeviceName string
}

type RefreshTokensInput struct {
	RefreshToken string
	UserAgent    string
	ClientIP     string
	DeviceName   string
}

type LogoutInput struct {
	RefreshToken string
}
