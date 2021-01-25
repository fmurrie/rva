package factory

import (
	// "context"
	"database/sql"
	"github.com/spf13/viper"
	"rva/dao"
	"rva/helper"
)

type RvaDao interface {
	DeployDatabase()
	OpenConnection() *sql.DB
	Execute(query string) interface{}
	ExecuteWithoutLock(query string) interface{}
	ExecuteContext(parameter interface{}, queries []string) interface{}
	LogError(err error)
}

func GetRvaDao() RvaDao {
	switch helper.GetRvaSecurityHelper().Decrypt(viper.GetString("database.driver")) {
	case "mysql":
		return dao.GetRvaMySqlDao()
	}
	return nil
}
