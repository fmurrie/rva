package model

import "gorm.io/gorm"

type AccountGroup struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;type:varchar(250);not null"`
	Accounts []Account `gorm:"many2many:role;"`
}