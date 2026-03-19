package models

import "time"

// Setting 系统设置
type Setting struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Category  string    `gorm:"size:50;index;not null" json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Setting) TableName() string {
	return "settings"
}
