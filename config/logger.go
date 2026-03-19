package config

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
}

// loadLoggerConfig 加载日志配置
func loadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:      getEnvOrDefault("LOG_LEVEL", "info"),
		Format:     getEnvOrDefault("LOG_FORMAT", "text"),
		Output:     getEnvOrDefault("LOG_OUTPUT", "both"),
		FilePath:   getEnvOrDefault("LOG_FILE", "./logs/app.log"),
		MaxSizeMB:  1,
		MaxBackups: 30,
		MaxAgeDays: 7,
		Compress:   true,
	}
}
