package nosql

import (
	"sync"
)

var instanceNoSqlDaoFactory *NoSqlDaoFactory
var singletonNoSqlDaoFactory sync.Once

type NoSqlDaoFactory struct {
}

func GetNoSqlDaoFactory() *NoSqlDaoFactory {
	singletonNoSqlDaoFactory.Do(func() {
		instanceNoSqlDaoFactory = &NoSqlDaoFactory{}
	})
	return instanceNoSqlDaoFactory
}

func (pointer NoSqlDaoFactory) GetNoSqlDao(key string) INoSqlDao {
	switch key {
	default:
		return nil
	}
}
