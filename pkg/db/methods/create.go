package methods

import (
	"gorm.io/gorm"
)

func Create[T any](pgdb *gorm.DB, checkPattern map[string]interface{}, obj *T) (status bool, err error) {
	var existing T
	err = pgdb.Where(checkPattern).First(&existing).Error
	if err == nil {
		return false, nil
	}

	if err != gorm.ErrRecordNotFound {
		return false, err
	}

	err = pgdb.Create(obj).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
