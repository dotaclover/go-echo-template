package config

import "time"

// AppConfig 应用配置
type AppConfig struct {
	Host             string
	Port             string
	Debug            bool
	Name             string
	Version          string
	Env              string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	ShutdownTimeout  time.Duration
	RateLimitEnabled bool
	RateLimitLimit   int
	RateLimitWindow  time.Duration
}

// loadAppConfig 加载应用配置
func loadAppConfig() AppConfig {
	return AppConfig{
		Host:             getEnvOrDefault("APP_HOST", "localhost"),
		Port:             getEnvOrDefault("APP_PORT", "8080"),
		Debug:            getEnvOrDefault("APP_DEBUG", "true") == "true",
		Name:             getEnvOrDefault("APP_NAME", "MyApp"),
		Version:          getEnvOrDefault("APP_VERSION", "0.1.0"),
		Env:              getEnvOrDefault("APP_ENV", "development"),
		ReadTimeout:      getDurationEnvOrDefault("APP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:     getDurationEnvOrDefault("APP_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:      getDurationEnvOrDefault("APP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:  getDurationEnvOrDefault("APP_SHUTDOWN_TIMEOUT", 30*time.Second),
		RateLimitEnabled: getEnvOrDefault("RATE_LIMIT_ENABLED", "false") == "true",
		RateLimitLimit:   getIntEnvOrDefault("RATE_LIMIT_LIMIT", 120),
		RateLimitWindow:  getDurationEnvOrDefault("RATE_LIMIT_WINDOW", time.Minute),
	}
}

// IsDebug 判断是否为调试模式
func IsDebug() bool {
	return getEnvOrDefault("APP_DEBUG", "true") == "true"
}
