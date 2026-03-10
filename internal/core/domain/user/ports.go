package user

import "time"

type IDGenerator interface {
	New() (string, error)
}

type Clock interface {
	Now() time.Time
}
