package helpers

import (
	"be-recipe/src/config"
	"be-recipe/src/models"
)

func Migrate() {
	config.DB.AutoMigrate(
		&models.User{},
		&models.Recipe{},
	)
}
