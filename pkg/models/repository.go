package models

type Sections struct {
	ID     uint   `gorm:"primaryKey"`
	Module string `gorm:"not null"`
	URL    string `gorm:"uniqueIndex;not null"`
}

type Permissions struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Sections    []Sections `gorm:"many2many:permissions_sections; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Roles struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Permissions []Permissions `gorm:"many2many:roles_permissions; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Users struct {
	ID       uint    `gorm:"primaryKey"`
	Username string  `gorm:"uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Status   bool    `gorm:"default:true"`
	Roles    []Roles `gorm:"many2many:users_roles; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
