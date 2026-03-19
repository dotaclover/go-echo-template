package config

import (
	"fmt"
	"time"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string
	MySQLDSN        string
	SQLitePath      string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	LogLevel        string
	SlowThreshold   time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// loadDatabaseConfig 加载数据库配置
func loadDatabaseConfig() DatabaseConfig {
	dbType := getEnvOrDefault("DB_TYPE", "sqlite")

	cfg := DatabaseConfig{
		Type:            dbType,
		MySQLDSN:        getEnvOrDefault("MYSQL_DSN", ""),
		SQLitePath:      getEnvOrDefault("SQLITE_PATH", "./data/app.db"),
		Host:            getEnvOrDefault("DB_HOST", "127.0.0.1"),
		Port:            getEnvOrDefault("DB_PORT", "3306"),
		User:            getEnvOrDefault("DB_USER", "root"),
		Password:        getEnvOrDefault("DB_PASSWORD", ""),
		Name:            getEnvOrDefault("DB_NAME", "myapp"),
		LogLevel:        getEnvOrDefault("DB_LOG_LEVEL", "warn"),
		SlowThreshold:   200 * time.Millisecond,
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
	}

	if dbType == "mysql" {
		cfg.MaxIdleConns = 10
		cfg.MaxOpenConns = 50
	}

	return cfg
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	switch c.Type {
	case "mysql":
		if c.MySQLDSN != "" {
			return c.MySQLDSN
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s",
			c.User, c.Password, c.Host, c.Port, c.Name)
	case "sqlite":
		return c.SQLitePath
	default:
		return ""
	}
}
