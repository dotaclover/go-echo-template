package commands

import (
	"fmt"
	"myapp/models"
	"myapp/server/database"
	"myapp/utils"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedCommand 填充测试数据命令
type SeedCommand struct{}

func (c *SeedCommand) Name() string        { return "seed" }
func (c *SeedCommand) Description() string { return "Seed database with test data" }

func (c *SeedCommand) Execute(args []string) error {
	_ = godotenv.Load()
	utils.InitLogger()

	utils.Logger.Info("Seeding database...")

	db := database.InitDB()
	defer database.CloseDB(db)

	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	utils.Logger.Info("Database seeding completed!")
	utils.Logger.Info("Test accounts:")
	utils.Logger.Info("  admin / admin123")
	utils.Logger.Info("  user01 / user123")
	return nil
}

func seedUsers(db *gorm.DB) error {
	users := []struct {
		Username string
		Password string
		Email    string
		RealName string
		Role     string
	}{
		{"admin", "admin123", "admin@example.com", "Admin", models.RoleAdmin},
		{"user01", "user123", "user01@example.com", "User", models.RoleUser},
	}

	for _, u := range users {
		var existing models.User
		if db.Where("username = ?", u.Username).First(&existing).Error == nil {
			utils.Logger.Infof("  - Skipped: %s (already exists)", u.Username)
			continue
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		user := models.User{
			Username:     u.Username,
			PasswordHash: string(hash),
			Email:        u.Email,
			RealName:     u.RealName,
			Role:         u.Role,
			Status:       models.StatusActive,
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
		utils.Logger.Infof("  + Created: %s (%s)", u.Username, u.Role)
	}
	return nil
}
