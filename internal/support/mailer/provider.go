package mailer

import "strings"

const (
	ProviderLog  = "log"
	ProviderSMTP = "smtp"
)

func normalizeProvider(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
