package pricehistory

import "time"

type PriceHistory struct {
	ProductID   string
	Price       string
	Currency    string
	RecordedAt  time.Time
	ChangeValue string
}
