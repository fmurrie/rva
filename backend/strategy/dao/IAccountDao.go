package dao

import "rva/model"

type IAccountDao interface {
	CreateAccount(account model.Account)
	GetAccountByFilter()
}