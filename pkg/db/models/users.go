package models

type Users struct {
	ID       uint    `gorm:"primaryKey"`
	Username string  `gorm:"uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Status   bool    `gorm:"default:true"`
	Roles    []Roles `gorm:"many2many:users_roles; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
