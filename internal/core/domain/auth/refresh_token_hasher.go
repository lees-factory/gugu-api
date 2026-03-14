package auth

type RefreshTokenHasher interface {
	Hash(value string) string
}
