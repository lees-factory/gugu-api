package auth

import "context"

type VerificationSender interface {
	SendVerification(ctx context.Context, email string, token string) error
}
