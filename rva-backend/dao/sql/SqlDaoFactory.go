package sql

import (
	"rva/keyword"
	"sync"
)

var instanceSqlDaoFactory *SqlDaoFactory
var singletonSqlDaoFactory sync.Once

type SqlDaoFactory struct {
}

func GetSqlDaoFactory() *SqlDaoFactory {
	singletonSqlDaoFactory.Do(func() {
		instanceSqlDaoFactory = &SqlDaoFactory{}
	})
	return instanceSqlDaoFactory
}

func (pointer SqlDaoFactory) GetSqlDao(database ISqlDatabase) ISqlDao {
	switch database.GetDriver() {
	default:
		return nil
	case keyword.MySql:
		return GetMySqlDao(database)
	}
}
