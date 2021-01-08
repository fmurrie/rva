package controller

import (
	"fmt"
	"rva/factory"
	"strings"
	"sync"
	"github.com/gin-gonic/gin"
)

var lockRvaRestController = &sync.Mutex{}
var instanceRvaRestController *RvaRestController

type RvaRestController struct {
	router *gin.Engine
	rvaDao factory.RvaDao
}

func GetRvaRestController() *RvaRestController {
	if instanceRvaRestController == nil {
		lockRvaRestController.Lock()
		defer lockRvaRestController.Unlock()
		if instanceRvaRestController == nil {
			instanceRvaRestController = &RvaRestController{
				router: gin.Default(),
				rvaDao: factory.GetRvaDao(),
			}
		}
	}
	return instanceRvaRestController
}

func (pointer RvaRestController) RaiseEndpoints() {

	result := pointer.rvaDao.ExecuteWithoutLock("call rvaEndpoint_getAll();")
	if result != nil {
		buildEndpoint := func(endpoint map[string]interface{}) {
			switch strings.ToUpper(endpoint["httpVerb"].(string)) {
			case "GET":
			case "HEAD":
			case "POST":
				pointer.router.POST(endpoint["path"].(string), func(context *gin.Context) {
					var json interface{}
					context.BindJSON(&json)
					result:=pointer.DoEndpointProcess(endpoint["idRvaEndpoint"].(int),json)
					context.JSON(200,result)
				})
			case "PUT":
			case "DELETE":
			case "CONNECT":
			case "OPTIONS":
			case "TRACE":
			case "PATCH":
			}
		}

		_, ok := result.(map[string]interface{})

		if ok {
			buildEndpoint(result.(map[string]interface{}))
		} else {
			for _, value := range result.([]map[string]interface{}) {
				buildEndpoint(value)
			}
		}

	}
	pointer.router.Run()
}

func (pointer RvaRestController) DoEndpointProcess(idRvaEndpoint int, parameter interface{}) interface{} {
	steps := pointer.rvaDao.ExecuteWithoutLock(fmt.Sprint("call rvaEndpointStep_getByIdRvaEndpoint(", idRvaEndpoint,");"))
	var queries []string
	if steps!=nil{
		switch steps.(type) {
		case map[string]interface{}:
			queries=append(queries, steps.(map[string]interface{})["procedureQuery"].(string))
		case []map[string]interface{}:
			for _, endpoint := range steps.([]map[string]interface{}) {
				queries=append(queries, endpoint["procedureQuery"].(string))
			}
		}
		result,err:=pointer.rvaDao.ExecuteContext(parameter,queries)
		if err!=nil{
			result=err
		}
		return result
	}
	return nil
}