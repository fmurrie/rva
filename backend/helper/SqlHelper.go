package helper

import (
	"sync"
)

var singletonSqlHelper *SqlHelper
var onceSqlHelper sync.Once

type SqlHelper struct {
}

func GetSqlHelper() *SqlHelper {
	onceSqlHelper.Do(func() {
		singletonSqlHelper = &SqlHelper{}
	})
	return singletonSqlHelper
}