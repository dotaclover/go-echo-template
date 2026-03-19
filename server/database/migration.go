package database

import (
	"myapp/models"
	"myapp/utils"

	"gorm.io/gorm"
)

// RunMigrations 执行数据库迁移
func RunMigrations(db *gorm.DB) error {
	utils.Logger.Info("Starting database migration...")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Setting{},
		// 在此添加新模型...
	); err != nil {
		return err
	}

	utils.Logger.Info("Database migration completed")
	return nil
}
