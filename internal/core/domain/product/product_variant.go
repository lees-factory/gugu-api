package product

import (
	"database/sql"
	"time"
)

type Variant struct {
	ProductID       string
	Language        string
	Currency        string
	Title           string
	MainImageURL    string
	ProductURL      string
	CurrentPrice    string
	LastCollectedAt *sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
