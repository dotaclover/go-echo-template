package services

import (
	"myapp/models"

	"gorm.io/gorm"
)

// SettingService 系统设置服务
type SettingService struct {
	db *gorm.DB
}

func NewSettingService(db *gorm.DB) *SettingService {
	return &SettingService{db: db}
}

// Get 获取设置值
func (s *SettingService) Get(key string) (string, error) {
	var setting models.Setting
	err := s.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

// Set 设置值（不存在则创建）
func (s *SettingService) Set(key, value, category string) error {
	var setting models.Setting
	err := s.db.Where("key = ?", key).First(&setting).Error

	if err == gorm.ErrRecordNotFound {
		return s.db.Create(&models.Setting{
			Key:      key,
			Value:    value,
			Category: category,
		}).Error
	}
	if err != nil {
		return err
	}

	setting.Value = value
	return s.db.Save(&setting).Error
}

// GetByCategory 获取某分类下的所有设置
func (s *SettingService) GetByCategory(category string) ([]models.Setting, error) {
	var settings []models.Setting
	err := s.db.Where("category = ?", category).Find(&settings).Error
	return settings, err
}
