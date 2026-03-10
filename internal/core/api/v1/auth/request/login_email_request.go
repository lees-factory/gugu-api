package request

type LoginEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
