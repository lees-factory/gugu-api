package user

import "time"

type Clock interface {
	Now() time.Time
}
