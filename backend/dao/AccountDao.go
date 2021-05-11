package dao

import (
	"rva/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AccountDao struct {
}

func (pointer AccountDao) CreateAccount(account *model.Account){
	dsn := "root:root@tcp(127.0.0.1:3306)/rva?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.Create(&account);
}

func (pointer AccountDao) GetAccountByEmailAndPassword() {
	dsn := "root:root@tcp(127.0.0.1:3306)/rva?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db:=db.Begin()
	db.b
}
