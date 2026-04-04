package product

import (
	"time"

	"github.com/ljj/gugu-api/internal/core/enum"
)

type Product struct {
	ID                string
	Market            enum.Market
	ExternalProductID string
	OriginalURL       string
	Title             string
	MainImageURL      string
	ProductURL        string
	PromotionLink     string
	CollectionSource  string
	LastCollectedAt   time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
