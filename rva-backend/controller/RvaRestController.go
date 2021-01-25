package controller

import (
	"net/http"
	"rva/factory"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var lockRvaRestController = &sync.Mutex{}
var instanceRvaRestController *RvaRestController

type RvaRestController struct {
	router                 *gin.Engine
	rvaDao                 factory.RvaDao
	rvaProcedureController *RvaProcedureController
	rvaAuthController      *RvaAuthController
}

func GetRvaRestController() *RvaRestController {
	if instanceRvaRestController == nil {
		lockRvaRestController.Lock()
		defer lockRvaRestController.Unlock()
		if instanceRvaRestController == nil {
			instanceRvaRestController = &RvaRestController{
				router:                 gin.Default(),
				rvaDao:                 factory.GetRvaDao(),
				rvaProcedureController: GetRvaProcedureController(),
				rvaAuthController:      GetRvaAuthController(),
			}
		}
	}
	return instanceRvaRestController
}

func (pointer RvaRestController) RaiseEndpoints() {
	result := pointer.rvaDao.ExecuteWithoutLock("call rvaEndpoint_getAll();")
	if result != nil {

		authValidation := func(valid bool, authHeader string) (bool, int, interface{}) {
			if valid {
				if authHeader == "" {
					return false, http.StatusBadRequest, map[string]interface{}{
						"Status":  false,
						"Message": "No token found",
					}
				}
				if pointer.rvaAuthController.ValidateToken(authHeader) {
					return true, http.StatusOK, nil
				}
				return false, http.StatusUnauthorized, map[string]interface{}{
					"Status":  false,
					"Message": "Token is not valid",
				}
			}
			return true, http.StatusOK, nil
		}

		authCreation := func(valid bool, account interface{}) {
			if _, ok := account.(map[string]interface{}); valid && ok {
				pointer.rvaAuthController.GenerateToken(account.(map[string]interface{}))
			}
		}

		buildEndpoint := func(endpoint map[string]interface{}) {
			switch strings.ToUpper(endpoint["httpVerb"].(string)) {
			case "GET":
			case "HEAD":
			case "POST":
				pointer.router.POST(endpoint["path"].(string), func(context *gin.Context) {
					var requestBody interface{}
					context.BindJSON(&requestBody)

					isValid, statusCode, response := authValidation(endpoint["validAuth"].(bool), context.GetHeader("Authorization"))
					if isValid || !endpoint["validAuth"].(bool) {
						response = pointer.rvaProcedureController.DoProcedure(endpoint["idRvaOrganizedProcedure"].(int), requestBody)
					}
					authCreation(endpoint["createAuth"].(bool), response)

					context.JSON(statusCode, response)
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
