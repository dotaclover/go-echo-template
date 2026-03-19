package config

import (
	"fmt"
	"os"
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

// IsValid 验证配置是否有效
func (c *Config) IsValid() error {
	if len(c.JWT.Secret) < 16 {
		return fmt.Errorf("JWT_SECRET must be at least 16 characters long")
	}
	if c.Database.Type == "" {
		return fmt.Errorf("DB_TYPE is required")
	}
	return nil
}
