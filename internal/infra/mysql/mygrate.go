package mysql

import (
	"hackfest-uc/internal/domain/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		entity.User{},
	); err != nil {
		return err
	}
	return nil
}
