package models

import "time"

// User 用户模型
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Email        string    `gorm:"type:varchar(100)" json:"email"`
	RealName     string    `gorm:"type:varchar(50)" json:"real_name"`
	Role         string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	Status       string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"

	StatusActive   = "active"
	StatusDisabled = "disabled"
)

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
