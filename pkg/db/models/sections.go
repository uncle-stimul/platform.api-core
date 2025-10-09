package models

type Sections struct {
	ID     uint   `gorm:"primaryKey"`
	Module string `gorm:"not null"`
	URL    string `gorm:"uniqueIndex;not null"`
}
