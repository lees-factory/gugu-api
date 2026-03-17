package mailer

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed templates/verification_email.html
var verificationEmailHTML string

type VerificationTemplateData struct {
	VerificationCode string
}

type VerificationTemplate struct {
	html *template.Template
}

func NewVerificationTemplate() (*VerificationTemplate, error) {
	parsed, err := template.New("verification_email").Parse(verificationEmailHTML)
	if err != nil {
		return nil, fmt.Errorf("parse verification email template: %w", err)
	}

	return &VerificationTemplate{html: parsed}, nil
}

func (t *VerificationTemplate) RenderHTML(data VerificationTemplateData) (string, error) {
	var buffer bytes.Buffer
	if err := t.html.Execute(&buffer, data); err != nil {
		return "", fmt.Errorf("render verification email template: %w", err)
	}

	return buffer.String(), nil
}
