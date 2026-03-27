package request

type Logout struct {
	RefreshToken string `json:"refresh_token"`
}
