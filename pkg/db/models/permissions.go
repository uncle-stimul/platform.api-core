package models

type Permissions struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Sections    []Sections `gorm:"many2many:permissions_sections; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
