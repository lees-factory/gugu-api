package email

import (
	"context"
	"fmt"
	"log"
)

type Sender interface {
	SendVerification(ctx context.Context, email string, code string) error
}

type Config struct {
	Provider     string
	MailFrom     string
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

type LogSender struct {
	template *VerificationTemplate
}

func NewSender(config Config) (Sender, error) {
	switch normalizeProvider(config.Provider) {
	case "", ProviderSMTP:
		return NewSMTPSender(
			config.MailFrom,
			config.SMTPHost,
			config.SMTPPort,
			config.SMTPUsername,
			config.SMTPPassword,
		)
	case ProviderLog:
		return NewLogSender()
	default:
		return nil, fmt.Errorf("unsupported mail provider: %s", config.Provider)
	}
}

func NewLogSender() (*LogSender, error) {
	template, err := NewVerificationTemplate()
	if err != nil {
		return nil, err
	}

	return &LogSender{template: template}, nil
}

func (s *LogSender) SendVerification(_ context.Context, email string, code string) error {
	html, err := s.template.RenderHTML(VerificationTemplateData{VerificationCode: code})
	if err != nil {
		return err
	}

	log.Printf("send verification email provider=log to=%s code=%s html=%q", email, code, html)
	return nil
}
