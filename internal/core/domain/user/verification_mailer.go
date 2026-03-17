package user

import "context"

type VerificationMailer interface {
	SendVerification(ctx context.Context, email string, code string) error
}
