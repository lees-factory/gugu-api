package config

import (
	"os"
	"strings"
)

type Config struct {
	HTTPAddress            string
	GoogleOAuthCallbackURL string
	AliExpressBaseURL      string
	AliExpressAppKey       string
	AliExpressAppSecret    string
	AliExpressCallbackURL  string
	DatabaseURL            string
	CORSAllowedOrigins     []string
	JWTSecret              string
	JWTIssuer              string
	MailProvider           string
	MailFrom               string
	SMTPHost               string
	SMTPPort               string
	SMTPUsername           string
	SMTPPassword           string
}

func Load() Config {
	return Config{
		HTTPAddress:            getenv("HTTP_ADDRESS", ":8080"),
		GoogleOAuthCallbackURL: getenv("GOOGLE_OAUTH_CALLBACK_URL", "http://localhost:8080/v1/auth/oauth/google/callback"),
		AliExpressBaseURL:      getenv("ALIEXPRESS_BASE_URL", "https://api-sg.aliexpress.com"),
		AliExpressAppKey:       getenv("ALIEXPRESS_APP_KEY", "528586"),
		AliExpressAppSecret:    os.Getenv("ALIEXPRESS_APP_SECRET"),
		AliExpressCallbackURL:  getenv("ALIEXPRESS_CALLBACK_URL", "https://googoo-client.vercel.app/callback"),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		CORSAllowedOrigins:     splitCSV(getenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		JWTSecret:              getenv("JWT_SECRET", "change-me"),
		JWTIssuer:              getenv("JWT_ISSUER", "gugu-api"),
		MailProvider:           getenv("MAIL_PROVIDER", "smtp"),
		MailFrom:               getenv("MAIL_FROM", "wjdrk70@gmail.com"),
		SMTPHost:               getenv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:               getenv("SMTP_PORT", "587"),
		SMTPUsername:           getenv("SMTP_USERNAME", "wjdrk70@gmail.com"),
		SMTPPassword:           os.Getenv("SMTP_PASSWORD"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func splitCSV(value string) []string {
	rawValues := strings.Split(value, ",")
	values := make([]string, 0, len(rawValues))

	for _, rawValue := range rawValues {
		trimmedValue := strings.TrimSpace(rawValue)
		if trimmedValue == "" {
			continue
		}
		values = append(values, trimmedValue)
	}

	return values
}
