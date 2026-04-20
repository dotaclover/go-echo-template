package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 应用配置结构
type Config struct {
	Database DatabaseConfig
	App      AppConfig
	JWT      JWTConfig
	Logger   LoggerConfig
}

// Load 加载配置
func Load() *Config {
	jwtConfig := loadJWTConfig()

	return &Config{
		Database: loadDatabaseConfig(),
		App:      loadAppConfig(),
		JWT:      jwtConfig,
		Logger:   loadLoggerConfig(),
	}
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// IsValid 验证配置是否有效
func (c *Config) IsValid() error {
	if len(c.JWT.Secret) < 16 {
		return fmt.Errorf("JWT_SECRET must be at least 16 characters long")
	}
	if c.JWT.RefreshExpiration <= 0 {
		return fmt.Errorf("JWT_REFRESH_EXPIRATION must be greater than 0")
	}
	if c.Database.Type == "" {
		return fmt.Errorf("DB_TYPE is required")
	}
	switch strings.ToLower(c.Database.Type) {
	case "sqlite":
		if c.Database.SQLitePath == "" {
			return fmt.Errorf("SQLITE_PATH is required when DB_TYPE=sqlite")
		}
	case "mysql":
		if c.Database.GetDSN() == "" {
			return fmt.Errorf("MYSQL_DSN or DB connection fields are required when DB_TYPE=mysql")
		}
	default:
		return fmt.Errorf("unsupported DB_TYPE: %s", c.Database.Type)
	}
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if c.App.ReadTimeout <= 0 || c.App.WriteTimeout <= 0 || c.App.IdleTimeout <= 0 || c.App.ShutdownTimeout <= 0 {
		return fmt.Errorf("app timeouts must be greater than 0")
	}
	if c.App.RateLimitEnabled && (c.App.RateLimitLimit <= 0 || c.App.RateLimitWindow <= 0) {
		return fmt.Errorf("rate limit configuration must be greater than 0")
	}
	return nil
}
