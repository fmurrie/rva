package helper

import (
	"github.com/spf13/viper"
	"sync"
)

var instanceConfigurationHelper = make(map[string]*ConfigurationHelper)
var singletonConfigurationHelper = make(map[interface{}]*sync.Once)

type ConfigurationHelper struct {
	fileName string
	path     string
	group    string
}

func GetConfigurationHelper(group string, fileName string, path string) *ConfigurationHelper {
	if singletonConfigurationHelper[path+"/"+fileName+"["+group+"]"]==nil{
		var once sync.Once
		singletonConfigurationHelper[path+"/"+fileName+"["+group+"]"]=&once
	}
	singletonConfigurationHelper[path+"/"+fileName+"["+group+"]"].Do(func() {
		instanceConfigurationHelper[path+"/"+fileName+"["+group+"]"] = &ConfigurationHelper{
			group:    group,
			fileName: fileName,
			path:     path,
		}
	})
	return instanceConfigurationHelper[path+"/"+fileName+"["+group+"]"]
}

func (pointer ConfigurationHelper) Read() {
	viper.AddConfigPath(pointer.path)
	viper.SetConfigName(pointer.fileName)
	viper.ReadInConfig()
}

func (pointer ConfigurationHelper) GetStringValueByKey(key string) string {
	pointer.Read()
	return viper.GetString(pointer.group + "." + key)
}

func (pointer ConfigurationHelper) GetBoolValueByKey(key string) bool {
	pointer.Read()
	return viper.GetBool(pointer.group + "." + key)
}

func (pointer ConfigurationHelper) GetIntValueByKey(key string) int {
	pointer.Read()
	return viper.GetInt(pointer.group + "." + key)
}

func (pointer ConfigurationHelper) GetFloat64ValueByKey(key string) float64 {
	pointer.Read()
	return viper.GetFloat64(pointer.group + "." + key)
}
