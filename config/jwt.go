package config

import (
	"os"
	"time"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

const DefaultJWTSecret = "change-me-to-a-secure-secret-key-at-least-32-chars"

// loadJWTConfig 加载JWT配置
func loadJWTConfig() JWTConfig {
	expiration := 24 * 7 * time.Hour // 7 days
	refreshExpiration := 24 * 14 * time.Hour
	if val := os.Getenv("JWT_EXPIRATION"); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			expiration = parsed
		}
	}
	if val := os.Getenv("JWT_REFRESH_EXPIRATION"); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			refreshExpiration = parsed
		}
	}

	return JWTConfig{
		Secret:            getEnvOrDefault("JWT_SECRET", DefaultJWTSecret),
		Expiration:        expiration,
		RefreshExpiration: refreshExpiration,
	}
}
