package factory

import (
	// "context"
	"database/sql"
	"rva/dao"
	"github.com/spf13/viper"
)


type RvaDao interface {
	DeployDatabase()
	OpenConnection() *sql.DB
	Execute(query string) interface{}
	ExecuteWithoutLock(query string) interface{}
	ExecuteContext(parameter interface{}, queries []string) (interface{}, error)
}

func GetRvaDao() RvaDao {
	switch viper.GetString("default.driver") {
	case "mysql":
		return dao.GetRvaMySqlDao()
	}
	return nil
}
