package auth

type IDGenerator interface {
	New() (string, error)
}
