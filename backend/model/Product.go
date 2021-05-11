package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name    string `json:"name" gorm:"unique;type:varchar(100);not null"`
}