package request

type VerifyEmail struct {
	Code string `json:"code"`
}
