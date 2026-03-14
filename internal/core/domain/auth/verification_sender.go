package auth

import "context"

type VerificationSender interface {
	SendVerification(ctx context.Context, email string, code string) error
}
