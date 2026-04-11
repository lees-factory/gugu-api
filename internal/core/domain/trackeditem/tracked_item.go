package trackeditem

import "time"

type TrackedItem struct {
	ID                    string
	UserID                string
	ProductID             string
	SKUID                 string
	OriginalURL           string
	ViewExternalProductID string
	PreferredLanguage     string
	TrackingScope         string
	Currency              string
	CreatedAt             time.Time
}
