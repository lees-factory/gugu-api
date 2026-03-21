package enum

import "strings"

type Market string

const (
	MarketAliExpress Market = "ALIEXPRESS"
	MarketCoupang    Market = "COUPANG"
	MarketEBay       Market = "EBAY"
)

func (m Market) Normalize() Market {
	return Market(strings.ToUpper(strings.TrimSpace(string(m))))
}

func (m Market) IsSupported() bool {
	switch m.Normalize() {
	case MarketAliExpress, MarketCoupang, MarketEBay:
		return true
	default:
		return false
	}
}
