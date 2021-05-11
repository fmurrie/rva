package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;type:varchar(100);not null"`
	Password string `json:"password" gorm:"type:varchar(2000);not null"`
	Groups []AccountGroup `gorm:"many2many:role;"`
}