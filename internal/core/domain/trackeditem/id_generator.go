package trackeditem

type IDGenerator interface {
	New() (string, error)
}
