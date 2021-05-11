package application

import (
	"net/http"
	"rva/client"
	"rva/controller"
	"rva/dao"
	"rva/helper"
	"rva/service"
	"sync"
	"time"
)

var multitonApplication = make(map[interface{}]*Application)
var onceApplication = make(map[interface{}]*sync.Once)

type Application struct {
	fileName string
	path     string
}

func GetApplication(fileName string,path string) *Application {
	var locker sync.Mutex
	locker.Lock()
    defer locker.Unlock()
	if onceApplication[path+"/"+fileName] == nil {
		onceApplication[path+"/"+fileName] = &sync.Once{}
	}
	onceApplication[path+"/"+fileName].Do(func() {
		multitonApplication[path+"/"+fileName] = &Application{
			fileName:fileName,
			path:path,
		}
	})
	return multitonApplication[path+"/"+fileName]
}

func (pointer Application) Run() error{
	sqlDatabase :=client.GetMySqlDBClient(helper.GetConfigurationHelper("database", pointer.fileName, pointer.path))
	//sqlDatabase.DeployDatabase()
	sqlDao:=dao.GetSqlDaoFactory().GetSqlDao(sqlDatabase)
	jwtService:=service.GetJwtService(client.GetJwtClient(helper.GetConfigurationHelper("jwt", pointer.fileName, pointer.path)))
	restController:=controller.GetRestController(pointer,sqlDao,jwtService)
	
	configurationHelper:=helper.GetConfigurationHelper("server", pointer.fileName, pointer.path)
	securityHelper:=helper.GetSecurityHelper()
	server := &http.Server{
        Addr:         ":"+securityHelper.Decrypt(configurationHelper.GetStringValueByKey("port")),
        Handler:      restController.GetRouter(),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

	return server.ListenAndServe()
}