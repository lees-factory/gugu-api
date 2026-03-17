package user

type CodeGenerator interface {
	New() (string, error)
}
