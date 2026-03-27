package pricesnapshot

import "time"

type SKUPriceSnapshot struct {
	SKUID         string
	SnapshotDate  time.Time
	Price         string
	OriginalPrice string
	Currency      string
}
