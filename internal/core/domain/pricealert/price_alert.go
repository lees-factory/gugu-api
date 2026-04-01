package pricealert

import "time"

type PriceAlert struct {
	ID        string
	UserID    string
	ProductID string
	Channel   string
	Enabled   bool
	CreatedAt time.Time
}
