package request

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	DeviceName   string `json:"device_name"`
}
