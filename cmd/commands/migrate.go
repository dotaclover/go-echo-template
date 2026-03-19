package commands

import (
	"fmt"
	"myapp/server/database"
	"myapp/utils"

	"github.com/joho/godotenv"
)

// MigrateCommand 数据库迁移命令
type MigrateCommand struct{}

func (c *MigrateCommand) Name() string        { return "migrate" }
func (c *MigrateCommand) Description() string { return "Run database migrations" }

func (c *MigrateCommand) Execute(args []string) error {
	_ = godotenv.Load()
	utils.InitLogger()

	utils.Logger.Info("Starting database migration...")

	db := database.InitDB()
	defer database.CloseDB(db)

	if err := database.RunMigrations(db); err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	utils.Logger.Info("Database migration completed successfully!")
	return nil
}
