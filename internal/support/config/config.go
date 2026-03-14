package config

import "os"

type Config struct {
	HTTPAddress            string
	GoogleOAuthCallbackURL string
	DatabaseURL            string
}

func Load() Config {
	return Config{
		HTTPAddress:            getenv("HTTP_ADDRESS", ":8080"),
		GoogleOAuthCallbackURL: getenv("GOOGLE_OAUTH_CALLBACK_URL", "http://localhost:8080/v1/auth/oauth/google/callback"),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
