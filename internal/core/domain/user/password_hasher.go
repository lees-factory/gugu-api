package user

type PasswordHasher interface {
	Hash(password string) (string, error)
}
