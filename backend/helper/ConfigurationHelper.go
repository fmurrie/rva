package helper

import (
	"sync"
	"github.com/spf13/viper"
)

var multitonConfigurationHelper = make(map[interface{}]*ConfigurationHelper)
var onceConfigurationHelper = make(map[interface{}]*sync.Once)

type ConfigurationHelper struct {
	reader *viper.Viper
	fileName string
	path     string
	group    string
}

func GetConfigurationHelper(group string, fileName string, path string) *ConfigurationHelper {
	if onceConfigurationHelper[path+"/"+fileName+"["+group+"]"]==nil{
		onceConfigurationHelper[path+"/"+fileName+"["+group+"]"]=&sync.Once{}
	}
	onceConfigurationHelper[path+"/"+fileName+"["+group+"]"].Do(func() {
		multitonConfigurationHelper[path+"/"+fileName+"["+group+"]"] = &ConfigurationHelper{
			reader: viper.New(),
			group:    group,
			fileName: fileName,
			path:     path,
		}
	})
	return multitonConfigurationHelper[path+"/"+fileName+"["+group+"]"]
}

func (pointer ConfigurationHelper) Read() {
	pointer.reader.AddConfigPath(pointer.path)
	pointer.reader.SetConfigName(pointer.fileName)
	pointer.reader.ReadInConfig()
}

func (pointer ConfigurationHelper) GetStringValueByKey(key string) string {
	pointer.Read()
	return pointer.reader.GetString(pointer.group + "." + key)
}

func (pointer ConfigurationHelper) GetBoolValueByKey(key string) bool {
	pointer.Read()
	return pointer.reader.GetBool(pointer.group + "." + key)
}

func (pointer ConfigurationHelper) GetIntValueByKey(key string) int {
	pointer.Read()
	return pointer.reader.GetInt(pointer.group + "." + key)
}

func (pointer ConfigurationHelper) GetFloat64ValueByKey(key string) float64 {
	pointer.Read()
	return pointer.reader.GetFloat64(pointer.group + "." + key)
}

func (pointer ConfigurationHelper) GetKey() string {
	return pointer.path + "/" + pointer.fileName + "[" + pointer.group + "]"
}
