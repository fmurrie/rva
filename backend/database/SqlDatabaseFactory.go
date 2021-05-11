package database

import (
	"sync"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var singletonSqlDatabaseFactory *SqlDatabaseFactory
var onceSqlDatabaseFactory sync.Once

type SqlDatabaseFactory struct {
}

func GetSqlDatabaseFactory() *SqlDatabaseFactory {
	onceSqlDatabaseFactory.Do(func() {
		singletonSqlDatabaseFactory = &SqlDatabaseFactory{}
	})
	return singletonSqlDatabaseFactory
}

func (pointer SqlDatabaseFactory) GetSqlDatabase() (*gorm.DB,error){
	return gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/rva?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
}