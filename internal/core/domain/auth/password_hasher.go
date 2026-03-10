package auth

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword string, rawPassword string) error
}
