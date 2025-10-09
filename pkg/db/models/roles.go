package models

type Roles struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Permissions []Permissions `gorm:"many2many:roles_permissions; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
