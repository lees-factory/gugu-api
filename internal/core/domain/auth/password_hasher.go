package auth

type PasswordVerifier interface {
	Verify(hashedPassword string, rawPassword string) error
}
