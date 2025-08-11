package migrations

import (
	"go-phonebook/models"

	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Contact{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Phone{}); err != nil {
		return err
	}
	return nil
}
