package request

type VerifyEmail struct {
	Token string `json:"token"`
}
