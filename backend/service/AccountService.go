package service

import (
	"rva/model"
	"rva/strategy/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


type AccountService struct {
	accountDao dao.IAccountDao
}

func GetAccountService(accountDao dao.IAccountDao) *AccountService{
	return &AccountService{
		accountDao:accountDao,
	}
}

func (pointer AccountService) CreateAccount(account model.Account) {
	dsn := "root:root@tcp(127.0.0.1:3306)/rva?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&account)
	  
		return nil
	  })
}

func (pointer AccountService) GetAccountByEmailAndPassword() {

}