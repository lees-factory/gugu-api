package product

type Market string

const (
	MarketAliExpress Market = "ALIEXPRESS"
	MarketCoupang    Market = "COUPANG"
	MarketEBay       Market = "EBAY"
)

func (m Market) IsSupported() bool {
	switch m {
	case MarketAliExpress, MarketCoupang, MarketEBay:
		return true
	default:
		return false
	}
}
