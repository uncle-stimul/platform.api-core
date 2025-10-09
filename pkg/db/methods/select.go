package methods

import "gorm.io/gorm"

func SelectAll[T any](db *gorm.DB, preloads ...string) ([]T, error) {
	query := db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	var records []T
	result := db.Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}
