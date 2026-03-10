package user

import "time"

type User struct {
	ID              string
	Email           string
	DisplayName     string
	PasswordHash    string
	AuthSource      string
	EmailVerified   bool
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
}
