package skupricehistory

import "time"

type SKUPriceHistory struct {
	SKUID       string
	Price       string
	Currency    string
	RecordedAt  time.Time
	ChangeValue string
}
