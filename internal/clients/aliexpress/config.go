package aliexpress

import "net/http"

const (
	defaultBaseURL   = "https://api-sg.aliexpress.com"
	defaultPartnerID = "gugu-api"
	signMethodSHA256 = "sha256"
)

type Config struct {
	BaseURL     string
	AppKey      string
	AppSecret   string
	CallbackURL string
	PartnerID   string
	HTTPClient  *http.Client
}
