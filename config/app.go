package config

// AppConfig 应用配置
type AppConfig struct {
	Host    string
	Port    string
	Debug   bool
	Name    string
	Version string
}

// loadAppConfig 加载应用配置
func loadAppConfig() AppConfig {
	return AppConfig{
		Host:    getEnvOrDefault("APP_HOST", "localhost"),
		Port:    getEnvOrDefault("APP_PORT", "8080"),
		Debug:   getEnvOrDefault("APP_DEBUG", "true") == "true",
		Name:    "MyApp",
		Version: "0.1.0",
	}
}

// IsDebug 判断是否为调试模式
func IsDebug() bool {
	return getEnvOrDefault("APP_DEBUG", "true") == "true"
}
