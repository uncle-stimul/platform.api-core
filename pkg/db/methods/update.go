package methods

import (
	"errors"

	"gorm.io/gorm"
)

func Update[T any](db *gorm.DB, id uint, obj *T) error {
	var existing T
	if err := db.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("объект для обновления не найден")
		}
		return err
	}

	return db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&existing).Updates(obj).Error
}
