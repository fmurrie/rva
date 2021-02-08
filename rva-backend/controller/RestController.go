package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rva/controller/auth"
	"rva/dao/sql"
	"rva/helper"
	"strings"
	"sync"
)

var instanceRestController *RestController
var singletonRestController sync.Once

type RestController struct {
	router         *gin.Engine
	sqlDao         sql.ISqlDao
	authController auth.IAuthController
	mapperHelper   *helper.MapperHelper
}

func GetRestController(sqlDao sql.ISqlDao, authController auth.IAuthController) *RestController {
	singletonRestController.Do(func() {
		instanceRestController = &RestController{
			router:         gin.Default(),
			sqlDao:         sqlDao,
			authController: authController,
			mapperHelper:   helper.GetMapperHelper(),
		}
	})
	return instanceRestController
}

func (pointer RestController) RaiseEndpoints() {
	buildEndpoint := func(endpoint map[string]interface{}) {
		pointer.router.Handle(strings.ToUpper(endpoint["httpVerb"].(string)), endpoint["path"].(string), func(context *gin.Context) {
			var request interface{}
			var statusCode int
			var response interface{}
			var isValid bool

			context.ShouldBindJSON(&request)

			if request == nil {
				request = make(map[string]interface{})
			}
			for _, value := range context.Params {
				pointer.mapperHelper.MapKeyValue(request, value.Key, value.Value)
			}
			for key, value := range context.Request.URL.Query() {
				if len(value) > 0 {
					pointer.mapperHelper.MapKeyValue(request, key, value[0])
				}
			}
			if endpoint["validAuth"].(bool) {
				if context.GetHeader("Authorization") == "" {
					isValid, statusCode, response = false, http.StatusBadRequest, map[string]interface{}{
						"Status":  false,
						"Message": "No token found",
					}
				}
				if pointer.authController.ValidateToken(context.GetHeader("Authorization")) {
					isValid = true
				}
				isValid, statusCode, response = false, http.StatusUnauthorized, map[string]interface{}{
					"Status":  false,
					"Message": "Token is not valid",
				}
			} else {
				isValid = true
			}

			if isValid || !endpoint["validAuth"].(bool) {
				response = pointer.sqlDao.ExecuteProcedureGroup(endpoint["idProcedureGroup"].(int), request)
			}

			if endpoint["createAuth"].(bool) && response != nil {
				context.Header("Authorization", pointer.authController.GenerateToken(response.(map[string]interface{})))
			}

			if statusCode == 0 {
				if response != nil {
					statusCode = http.StatusOK
				} else {
					statusCode = http.StatusNoContent
				}
			}

			if statusCode == 0 {
				if response != nil {
					statusCode = endpoint["statusCodeSuccess"].(int)
				} else {
					statusCode = endpoint["statusCodeNoResponse"].(int)
				}
			}
			if endpoint["retriveOnlyHeaders"].(bool) {
				function := func(parameter map[string]interface{}) {
					for key, value := range parameter {
						context.Header(key, fmt.Sprint(value))
					}
				}
				pointer.mapperHelper.ParameterizedRecursionNotReturned(function, response)
				context.Status(statusCode)
			} else {
				context.JSON(statusCode, response)
			}
		})
	}

	pointer.mapperHelper.ParameterizedRecursionNotReturned(buildEndpoint, pointer.sqlDao.GetEndpoints())
	pointer.router.Run()
}
