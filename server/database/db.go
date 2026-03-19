package database

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/utils"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	gosqlite "github.com/glebarez/sqlite"
)

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	cfg := config.Load()
	return InitDBWithConfig(cfg.Database)
}

// InitDBWithConfig 使用配置初始化数据库连接
func InitDBWithConfig(cfg config.DatabaseConfig) *gorm.DB {
	var db *gorm.DB
	var err error

	// GORM 日志
	var level logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
	default:
		level = logger.Warn
	}

	gormLogger := utils.NewGormLogrusLogger(logger.Config{
		LogLevel:                  level,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
		SlowThreshold:             cfg.SlowThreshold,
	})

	gormConfig := &gorm.Config{Logger: gormLogger}

	dsn := cfg.GetDSN()
	if dsn == "" {
		utils.Logger.Fatalf("Database DSN is empty for type: %s", cfg.Type)
	}

	switch cfg.Type {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	case "sqlite":
		dir := filepath.Dir(dsn)
		if dir != "" && dir != "." {
			if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
				utils.Logger.Fatalf("Failed to create SQLite directory %s: %v", dir, mkErr)
			}
		}
		db, err = gorm.Open(gosqlite.Open(dsn), gormConfig)
	default:
		utils.Logger.Fatalf("Unsupported database type: %s (supported: mysql, sqlite)", cfg.Type)
	}

	if err != nil {
		utils.Logger.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		utils.Logger.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	utils.Logger.Infof("Database connected: %s", cfg.Type)
	return db
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		utils.Logger.Errorf("Failed to get sql.DB for close: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		utils.Logger.Errorf("Failed to close database: %v", err)
	} else {
		utils.Logger.Info("Database connection closed")
	}
}

// HealthCheck 数据库健康检查
func HealthCheck(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}
