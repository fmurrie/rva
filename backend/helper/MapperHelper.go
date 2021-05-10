package helper

import (
	"sync"
)

var singletonMapperHelper *MapperHelper
var onceMapperHelper sync.Once

type MapperHelper struct {
}

func GetMapperHelper() *MapperHelper {
	onceMapperHelper.Do(func() {
		singletonMapperHelper = &MapperHelper{}
	})
	return singletonMapperHelper
}