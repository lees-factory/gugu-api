package response

import (
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/user"
)

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	DisplayName   string    `json:"display_name"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewUser(source user.User) User {
	return User{
		ID:            source.ID,
		Email:         source.Email,
		DisplayName:   source.DisplayName,
		EmailVerified: source.EmailVerified,
		CreatedAt:     source.CreatedAt,
	}
}
