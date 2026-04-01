package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

const verificationSubject = "Gugu Plus 이메일 인증 코드"

type SMTPSender struct {
	template *VerificationTemplate
	from     string
	host     string
	port     string
	username string
	password string
}

func NewSMTPSender(from string, host string, port string, username string, password string) (*SMTPSender, error) {
	template, err := NewVerificationTemplate()
	if err != nil {
		return nil, err
	}
	if from == "" {
		return nil, fmt.Errorf("mail from is empty")
	}
	if host == "" || port == "" || username == "" || password == "" {
		return nil, fmt.Errorf("smtp configuration is incomplete")
	}

	return &SMTPSender{
		template: template,
		from:     from,
		host:     host,
		port:     port,
		username: username,
		password: password,
	}, nil
}

func (s *SMTPSender) SendVerification(_ context.Context, email string, code string) error {
	html, err := s.template.RenderHTML(VerificationTemplateData{VerificationCode: code})
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	message := buildHTMLMessage(s.from, email, verificationSubject, html)
	address := s.host + ":" + s.port

	if err := smtp.SendMail(address, auth, s.from, []string{email}, []byte(message)); err != nil {
		return fmt.Errorf("send smtp mail to %s: %w", email, err)
	}

	return nil
}

const priceAlertSubject = "Gugu Plus 가격 변동 알림"

func (s *SMTPSender) SendPriceAlert(_ context.Context, email string, productTitle string, oldPrice string, newPrice string, currency string) error {
	html := buildPriceAlertHTML(productTitle, oldPrice, newPrice, currency)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	message := buildHTMLMessage(s.from, email, priceAlertSubject, html)
	address := s.host + ":" + s.port

	if err := smtp.SendMail(address, auth, s.from, []string{email}, []byte(message)); err != nil {
		return fmt.Errorf("send price alert mail to %s: %w", email, err)
	}
	return nil
}

func buildPriceAlertHTML(productTitle string, oldPrice string, newPrice string, currency string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><body style="font-family:sans-serif;padding:20px;">
<h2>가격 변동 알림</h2>
<p><strong>%s</strong></p>
<p>이전 가격: %s %s</p>
<p>현재 가격: <strong>%s %s</strong></p>
<p style="color:#888;font-size:12px;">Gugu Plus에서 발송된 알림입니다.</p>
</body></html>`, productTitle, oldPrice, currency, newPrice, currency)
}

func buildHTMLMessage(from string, to string, subject string, html string) string {
	headers := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + encodeHeader(subject),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + html
}

func encodeHeader(value string) string {
	return "=?UTF-8?B?" + base64Encode(value) + "?="
}
