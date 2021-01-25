package controller

import (
	"fmt"
	"rva/factory"
	"sync"
)

var lockRvaProcedureController = &sync.Mutex{}
var instanceRvaProcedureController *RvaProcedureController

type RvaProcedureController struct {
	rvaDao factory.RvaDao
}

func GetRvaProcedureController() *RvaProcedureController {
	if instanceRvaProcedureController == nil {
		lockRvaProcedureController.Lock()
		defer lockRvaProcedureController.Unlock()
		if instanceRvaProcedureController == nil {
			instanceRvaProcedureController = &RvaProcedureController{
				rvaDao: factory.GetRvaDao(),
			}
		}
	}
	return instanceRvaProcedureController
}

func (pointer RvaProcedureController) DoProcedure(idRvaOrganizedProcedure int, parameter interface{}) interface{} {
	steps := pointer.rvaDao.ExecuteWithoutLock(fmt.Sprint("call rvaOrganizedProcedureStep_getByIdRvaOrganizedProcedure(", idRvaOrganizedProcedure, ");"))
	var queries []string
	if steps != nil {
		switch steps.(type) {
		case map[string]interface{}:
			queries = append(queries, steps.(map[string]interface{})["procedureQuery"].(string))
		case []map[string]interface{}:
			for _, endpoint := range steps.([]map[string]interface{}) {
				queries = append(queries, endpoint["procedureQuery"].(string))
			}
		}
		return pointer.rvaDao.ExecuteContext(parameter, queries)
	}
	return nil
}
