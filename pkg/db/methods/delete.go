package methods

import (
	"errors"

	"gorm.io/gorm"
)

func Delete[T any](pgdb *gorm.DB, id uint) error {
	var obj T
	err := pgdb.First(&obj, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("объект для удаления не найден")
		}
		return err
	}

	return pgdb.Delete(&obj).Error
}
