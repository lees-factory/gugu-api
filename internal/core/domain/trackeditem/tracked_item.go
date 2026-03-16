package trackeditem

import "time"

type TrackedItem struct {
	ID          string
	UserID      string
	ProductID   string
	OriginalURL string
	CreatedAt   time.Time
}
