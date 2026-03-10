package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/user"
)

type UserEmailVerifier interface {
	Verify(ctx context.Context, userID string, verifiedAt time.Time) error
}

type userEmailVerifier struct {
	userWriter user.Writer
}

func NewUserEmailVerifier(userWriter user.Writer) UserEmailVerifier {
	return &userEmailVerifier{userWriter: userWriter}
}

func (v *userEmailVerifier) Verify(ctx context.Context, userID string, verifiedAt time.Time) error {
	if err := v.userWriter.MarkEmailVerified(ctx, userID, verifiedAt); err != nil {
		return fmt.Errorf("mark email verified: %w", err)
	}
	return nil
}
