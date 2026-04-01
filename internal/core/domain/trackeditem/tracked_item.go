package trackeditem

import "time"

type TrackedItem struct {
	ID          string
	UserID      string
	ProductID   string
	SKUID       string
	OriginalURL string
	Currency    string
	CreatedAt   time.Time
}
