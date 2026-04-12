package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	AppEnv                string
	Port                  string
	DatabaseURL           string
	RedisURL              string
	RedisPassword         string
	GeminiAPIKey          string
	GeminiModel           string
	GeminiFallbackModels  []string
	JWTSecret             string
	JWTAccessTTL          time.Duration
	JWTRefreshTTL         time.Duration
	AdminUsername         string
	AdminPassword         string
	AllowedOrigins        []string
	TelegramToken         string
	TelegramChatID        string
	OrderLimitEnabled     bool
	OutboundProxyScheme   string
	OutboundProxyHost     string
	OutboundProxyPort     string
	OutboundProxyUsername string
	OutboundProxyPassword string
}

func Load() *Config {
	return &Config{
		AppEnv:                getEnv("ENV", "development"),
		Port:                  getEnv("PORT", "8080"),
		DatabaseURL:           getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/rostmebel?sslmode=disable"),
		RedisURL:              getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		GeminiAPIKey:          getEnv("GEMINI_API_KEY", ""),
		GeminiModel:           getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		GeminiFallbackModels:  getListEnv("GEMINI_FALLBACK_MODELS", "gemini-2.5-flash-lite"),
		JWTSecret:             getEnv("JWT_SECRET", "default-secret"),
		JWTAccessTTL:          getDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL:         getDurationEnv("JWT_REFRESH_TTL", 720*time.Hour),
		AdminUsername:         getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:         getEnv("ADMIN_PASSWORD", "admin"),
		AllowedOrigins:        getListEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:80"),
		TelegramToken:         getEnv("TELEGRAM_TOKEN", ""),
		TelegramChatID:        getEnv("TELEGRAM_CHAT_ID", ""),
		OrderLimitEnabled:     getEnv("ORDER_LIMIT_ENABLED", "true") == "true",
		OutboundProxyScheme:   getEnv("OUTBOUND_PROXY_SCHEME", "http"),
		OutboundProxyHost:     getEnv("OUTBOUND_PROXY_HOST", ""),
		OutboundProxyPort:     getEnv("OUTBOUND_PROXY_PORT", ""),
		OutboundProxyUsername: getEnv("OUTBOUND_PROXY_USERNAME", ""),
		OutboundProxyPassword: getEnv("OUTBOUND_PROXY_PASSWORD", ""),
	}
}

func getListEnv(key, fallback string) []string {
	raw := getEnv(key, fallback)
	items := strings.Split(raw, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(value)
		if err == nil {
			return d
		}
	}
	return fallback
}
