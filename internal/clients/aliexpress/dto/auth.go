package dto

type TokenExchangeInput struct {
	Code string
}

type RefreshTokenInput struct {
	RefreshToken string
}

type TokenSet struct {
	AccessToken           string
	RefreshToken          string
	ExpiresIn             int64
	RefreshExpiresIn      int64
	ExpireTime            int64
	RefreshTokenValidTime int64
	HavanaID              string
	UserID                string
	SellerID              string
	UserNick              string
	Account               string
	Locale                string
	AccountPlatform       string
	SP                    string
	RequestID             string
	Code                  string
}
