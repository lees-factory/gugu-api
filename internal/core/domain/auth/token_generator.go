package auth

type TokenGenerator interface {
	New() (string, error)
}
