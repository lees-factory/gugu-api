package pricesnapshot

import "time"

type ProductPriceSnapshot struct {
	ProductID    string
	SnapshotDate time.Time
	Price        string
	Currency     string
}
