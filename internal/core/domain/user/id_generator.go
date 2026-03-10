package user

type IDGenerator interface {
	New() (string, error)
}
