package methods

import (
	"errors"

	"gorm.io/gorm"
)

func SetAssociations[T any, A any](pgdb *gorm.DB, obj *T, assocName string, assoc []A) error {
	var existing T
	err := pgdb.First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("основной объект не существует")
		}
		return err
	}

	return pgdb.Model(obj).Association(assocName).Append(assoc)
}
