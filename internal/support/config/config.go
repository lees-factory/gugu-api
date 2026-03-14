package config

import "os"

type Config struct {
	HTTPAddress            string
	GoogleOAuthCallbackURL string
	DatabaseURL            string
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
		DatabaseURL:            os.Getenv("DATABASE_URL"),
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
