package user

import "time"

type NewUser struct {
	Email           string
	DisplayName     string
	PasswordHash    string
	AuthSource      string
	EmailVerified   bool
	EmailVerifiedAt *time.Time
}
