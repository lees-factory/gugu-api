package email

import (
	"context"
	"log"
)

type LogSender struct{}

func (LogSender) SendVerification(_ context.Context, email string, token string) error {
	log.Printf("send verification email to=%s token=%s", email, token)
	return nil
}
